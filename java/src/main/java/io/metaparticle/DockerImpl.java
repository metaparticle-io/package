package io.metaparticle;

import java.io.OutputStream;
import java.util.ArrayList;
import io.metaparticle.annotations.Runtime;
import static io.metaparticle.Util.handleErrorExec;
import static io.metaparticle.Util.addAll;

public class DockerImpl implements Executor, Builder {
    @Override
    public boolean build(String dir, String image, OutputStream stdout, OutputStream stderr) {
        return handleErrorExec(new String[] {"docker", "build", "-t", image, dir}, stdout, stderr);
    }

    @Override
    public boolean push(String image, OutputStream stdout, OutputStream stderr) {
        return handleErrorExec(new String[] {"docker", "push", image}, stdout, stderr);
    }

    @Override
    public boolean run(String image, String name, Runtime config, OutputStream stdout, OutputStream stderr) {
        ArrayList<String> cmd = new ArrayList<>();
        addAll(cmd, new String[]{"docker", "run", "-d", "--name", name});
        if (config != null && config.ports() != null) {
            for (int port : config.ports()) {
                cmd.add("-p=" + port + ":" + port);
            }
        }
        cmd.add(image);
        return handleErrorExec(cmd.toArray(new String[]{}),  stdout, stderr);
    }

    @Override
    public boolean logs(String name, OutputStream stdout, OutputStream stderr) {
        return handleErrorExec(new String[] {"docker", "logs", "-f", name}, stdout, stderr);
    }

    @Override
    public boolean cancel(String name) {
        return handleErrorExec(new String[] {"docker", "rm", "-f", name}, null, null);
    }
}
