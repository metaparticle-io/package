package io.metaparticle;

import java.io.File;
import java.io.IOException;
import java.io.OutputStream;
import java.lang.reflect.Array;
import java.lang.reflect.Method;
import java.nio.file.Files;

import io.metaparticle.annotations.Package;
import io.metaparticle.annotations.Runtime;
import static io.metaparticle.Util.handleErrorExec;
import static io.metaparticle.Util.once;

public class Metaparticle {

    public static boolean inDockerContainer() {
        if ("true".equals(System.getenv("METAPARTICLE_IN_CONTAINER"))) {
            return true;
        }
        File f = new File("/proc/1/cgroup");
        if (f.exists()) {
            try {
                String s = new String(Files.readAllBytes(f.toPath()), "UTF-8");
                return s.indexOf("docker") != -1;
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
            
            String name = "web";
            String image = "test";

            try {
                Class clazz = Class.forName(className);
                Method m = clazz.getMethod(methodName, String[].class);
                Package p = m.getAnnotation(Package.class);
                Runtime r = m.getAnnotation(Runtime.class);

                Executor exec = getExecutor(r);
                Builder builder = getBuilder(p);

                OutputStream stdout = p.verbose() ? System.out : null;
                OutputStream stderr = p.quiet() ? null : System.err;

                //handleErrorExec(new String[] {"mvn", "package"}, System.out, System.err);

                builder.build(".", image, stdout, stderr);

                Runnable cancel = once(() -> exec.cancel(name));

                java.lang.Runtime.getRuntime().addShutdownHook(new Thread(cancel));

                exec.run(image, name, r, stdout, stderr);
                exec.logs(name, System.out, System.err);
                cancel.run();
            } catch (NoSuchMethodException | ClassNotFoundException ex) {
                // This should really never happen.
                throw new IllegalStateException(ex);
            }
        }
    }
}