package io.metaparticle;

import static io.metaparticle.Util.handleErrorExec;
import static io.metaparticle.Util.once;

import java.io.File;
import java.io.IOException;
import java.io.OutputStream;
import java.lang.reflect.Method;
import java.nio.file.Files;
import java.nio.file.Paths;

import io.metaparticle.annotations.Package;
import io.metaparticle.annotations.Runtime;

public class Metaparticle {

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
                return (s.indexOf("docker") != -1 || s.indexOf("kubepods") != -1);
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

    public static void writeDockerfile(String className, Package p) throws IOException {
        byte [] output;
        if (p.dockerfile() == null || p.dockerfile().length() == 0) {
            String contents = 
"FROM openjdk:8-jre-alpine\n" +
"COPY %s /main.jar\n" +
"CMD java -classpath /main.jar %s";
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

                writeDockerfile(className, packageAnnotation);

                if (packageAnnotation.build()) {
                    handleErrorExec(new String[] {"mvn", "package"}, System.out, System.err);                    
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
}
