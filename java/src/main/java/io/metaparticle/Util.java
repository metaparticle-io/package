package io.metaparticle;

import java.io.InputStream;
import java.io.IOException;
import java.io.OutputStream;

import com.google.common.io.ByteStreams;

public class Util {
    protected static <T> void addAll(java.util.List<T> list, T[] arr) {
        for(T elt : arr) {
            list.add(elt);
        }
    }

    private static void safeCopy(InputStream is, OutputStream os) {
        try {
            ByteStreams.copy(is, os);
        } catch (IOException ex) {
            ex.printStackTrace();
        }
    }

    protected static boolean handleErrorExec(String[] args, OutputStream stdout, OutputStream stderr) {
        try {
            Process proc = Runtime.getRuntime().exec(args);
            if (stdout != null) {
                new Thread(() -> safeCopy(proc.getInputStream(), stdout)).start();
            }
            if (stderr != null) {
                new Thread(() -> safeCopy(proc.getErrorStream(), stderr)).start();
            }
            int code = proc.waitFor();
            if (code != 0 && stderr == null) {
                ByteStreams.copy(proc.getErrorStream(), System.err);
            }
            return code == 0;
        } catch (IOException | InterruptedException ex) {
            ex.printStackTrace();
            return false;
        }
    }

    protected static Runnable once(final Runnable r) {
        return new Runnable() {
            boolean run = false;
            public void run() {
                synchronized (this) {
                    if (!run) {
                        r.run();
                        run = true;
                    }
                }
            }
        };
    }
}