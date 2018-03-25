package io.metaparticle;

import io.metaparticle.annotations.Package;
import io.metaparticle.annotations.Runtime;
import java.io.File;
import java.io.IOException;
import java.io.OutputStream;
import java.lang.reflect.Method;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.logging.Logger;

import static io.metaparticle.Util.handleErrorExec;
import static io.metaparticle.Util.once;

public class Metaparticle {

    private static Logger logger = Logger.getLogger(Metaparticle.class.getName());

    static final String SPRING_BOOT_APP_ANNOTATION = "org.springframework.boot.autoconfigure.SpringBootApplication";

    static final String SPRING_BOOT_LOADER_JAR_LAUNCHER_CLASS = "org.springframework.boot.loader.JarLauncher";

    public static boolean inDockerContainer() {
        String envFlag = System.getenv("METAPARTICLE_IN_CONTAINER");
        if (envFlag != null) {
            switch (envFlag) {
                case "true":
                case "1":
                    return true;
                case "false":
                case "0":
                    return false;
            }
        }

        File f = new File("/proc/1/cgroup");
        if (f.exists()) {
            try {
                String s = new String(Files.readAllBytes(f.toPath()), "UTF-8");
                return (s.contains("docker") || s.contains("kubepods"));
            } catch (IOException ex) {
                throw new IllegalStateException(ex);
            }
        }
        return false;
    }

    public static Executor getExecutor(Runtime cfg) {
        if (cfg == null) {
            return new DockerImpl();
        }
        switch (cfg.executor()) {
            case "docker":
                return new DockerImpl();
            case "aci":
                return new AciExecutor();
            case "metaparticle":
                return new MetaparticleExecutor();
            default:
                throw new IllegalStateException("Unknown executor: " + cfg.executor());
        }
    }

    public static Builder getBuilder(Package pkg) {
        if (pkg == null) {
            return new DockerImpl();
        }
        switch (pkg.builder()) {
            case "docker":
                return new DockerImpl();
            default:
                throw new IllegalStateException("Unknown builder: " + pkg.builder());
        }
    }

    public static void writeDockerfile(String className, Package p, boolean runSpringBootJar) throws IOException {
        byte [] output;
        if (p.dockerfile().isEmpty()) {
            String contents =
"FROM openjdk:8-jre-alpine\n" +
"COPY %s /main.jar\n" +
(runSpringBootJar ? "CMD java -jar /main.jar" : "CMD java -classpath /main.jar %s");
            output = String.format(contents, p.jarFile(), className).getBytes();
        } else {
            output = Files.readAllBytes(Paths.get(p.dockerfile()));
        }
        Files.write(Paths.get("Dockerfile"), output);
    }

    public static void Containerize(Runnable fn) {
        if (inDockerContainer()) {
            fn.run();
        } else {
            File f = new File("pom.xml");
            try {
                if (!f.exists()) {
                    throw new IllegalStateException("Can not find: " + f.getCanonicalPath());
                }
            } catch (IOException ex) {
                throw new IllegalStateException(ex);
            }

            StackTraceElement[] traces = Thread.currentThread().getStackTrace();
            String className = traces[2].getClassName();
            String methodName = traces[2].getMethodName();

            try {
                Class<?> clazz = Class.forName(className);
                String name = clazz.getCanonicalName().replace('.', '-').toLowerCase();
                String image = name;
                Method method = clazz.getMethod(methodName, String[].class);
                Package packageAnnotation = method.getAnnotation(Package.class);
                Runtime runtimeAnnotation = method.getAnnotation(Runtime.class);

                if (packageAnnotation.repository().length() != 0) {
                    image = packageAnnotation.repository() + "/" + image;
                }

                Executor exec = getExecutor(runtimeAnnotation);
                Builder builder = getBuilder(packageAnnotation);

                OutputStream stdout = packageAnnotation.verbose() ? System.out : null;
                OutputStream stderr = packageAnnotation.quiet() ? null : System.err;

                writeDockerfile(className, packageAnnotation, isSpringBootAppBeingContainerized(clazz));

                if (packageAnnotation.build()) {
                    doPackage(isRunningAsSpringBootApplication(traces), targetJarExists(packageAnnotation.jarFile()));
                    builder.build(".", image, stdout, stderr);
                    if (packageAnnotation.publish()) {
                        builder.push(image, stdout, stderr);
                    }
                }

                Runnable cancel = once(() -> exec.cancel(name));
                java.lang.Runtime.getRuntime().addShutdownHook(new Thread(cancel));

                exec.run(image, name, runtimeAnnotation, stdout, stderr);
                exec.logs(name, System.out, System.err);
                cancel.run();
            } catch (NoSuchMethodException | ClassNotFoundException | IOException ex) {
                // This should really never happen.
                throw new IllegalStateException(ex);
            }
        }
    }

    static boolean isSpringBootAppBeingContainerized(Class<?> containerizedClass) {
        return Arrays.stream(containerizedClass.getAnnotations())
                .anyMatch(annotation -> annotation.toString().contains(SPRING_BOOT_APP_ANNOTATION));
    }

    static boolean targetJarExists(String jarfileName) {
        return Files.exists(Paths.get(jarfileName));
    }

    static boolean isRunningAsSpringBootApplication(StackTraceElement[] traces) {
        String launchClass = traces[traces.length - 1].getClassName();
        String launchMethod = traces[traces.length - 1].getMethodName();
        return launchClass.equals(SPRING_BOOT_LOADER_JAR_LAUNCHER_CLASS) && launchMethod.equals("main");
    }

    static boolean doPackage(boolean inRunningSpringBootApp, boolean targetExists) {
        if (inRunningSpringBootApp && targetExists) {
            logger.info("Package step skipped when running as Spring Boot Jar and target already exists");
            return false;
        } else {
            return handleErrorExec(new String[]{"mvn", "package"}, System.out, System.err);
        }
    }
}
