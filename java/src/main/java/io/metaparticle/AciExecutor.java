package io.metaparticle;

import java.io.OutputStream;
import java.util.ArrayList;

import static io.metaparticle.Util.addAll;
import static io.metaparticle.Util.handleErrorExec;

public class AciExecutor implements Executor {
    @Override
    public boolean cancel(String id) {
        return handleErrorExec(new String[] {"az", "container", "delete",  "-g", "test", "-n", id, "--yes"}, null, null);
    }

    @Override
    public boolean logs(String id, OutputStream stdout, OutputStream stderr) {
        while (true) {
            handleErrorExec(new String[] {"az", "container", "logs", "-g", "test", "-n", id}, stdout, stderr);
            try {
                Thread.sleep(5 * 1000);
            } catch (InterruptedException ex) {}
        }
    }

    @Override
    public boolean run(String image, String name, io.metaparticle.annotations.Runtime config, OutputStream stdout, OutputStream stderr) {
        String rg = "test";

        ArrayList<String> cmd = new ArrayList<>();
        addAll(cmd, new String[] {"az", "container", "create", "--image", image, "-g", rg, "-n", name, "--env=METAPARTICLE_IN_CONTAINER=true"});
        if (config != null && config.ports() != null && config.ports().length > 0) {
            cmd.add("--port=" + config.ports()[0]);
        }

        if (config != null && config.publicAddress()) {
            cmd.add("--ip-address=Public");
        }

        return handleErrorExec(cmd.toArray(new String[]{}), stdout, stderr);
    }
}