using System;
using System.Threading.Tasks;
using System.IO;
using System.Diagnostics;

namespace Metaparticle.Package
{
    public class Util
    {
        public static Process Exec(String file, String args, TextWriter stdout=null, TextWriter stderr=null)
        {
            var task = ExecAsync(file, args, stdout, stderr);
            task.Wait();
            return task.Result;
        }

        public static async Task Copy(StreamReader reader, TextWriter writer)
        {
            if (reader == null || writer == null) {
                return;
            }
            var read = await reader.ReadLineAsync();
            while (read != null) {
                await writer.WriteLineAsync(read);
                await writer.FlushAsync();
                read = await reader.ReadLineAsync();
            }
        }

        public static async Task<Process> ExecAsync(String file, String args, TextWriter stdout=null, TextWriter stderr=null)
        {
            var startInfo = new ProcessStartInfo(file, args);
            startInfo.RedirectStandardOutput = true;
            startInfo.RedirectStandardError = true;
            var proc = Process.Start(startInfo);
            var runTask = Task.Run(() => { proc.WaitForExit(); });
            var outputTask = Copy(proc.StandardOutput, stdout);
            var errTask = Copy(proc.StandardError, stderr);

            await Task.Run(() => 
            {
                Task.WaitAll(new []
                {
                    runTask,
                    outputTask,
                    errTask,
                });
            });
            return proc;
        }
    }
}