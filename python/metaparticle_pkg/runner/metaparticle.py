import json
import os
import subprocess

class MetaparticleRunner:
    def cancel(self, name):
        subprocess.check_call(['mp-compiler', '-f', '.metaparticle/spec.json', '--delete'])

    def logs(self, name):
        subprocess.check_call(['mp-compiler', '-f', '.metaparticle/spec.json', '--deploy=false', '--attach=true'])

    def ports(self, portArray):
        result = []
        for port in portArray:
            result.append({
                'number': port,
                'protocol': 'TCP'
            })
        return result

    def run(self, img, name, options):
        svc =  {
            "name": name,
            "guid": 1234567, 
            "services": [ 
                {
                    "name": name,
                    "replicas": options.replicas,
                    "shardSpec": options.shardSpec,
                    "containers": [
                        { "image": img }
                    ],
                    "ports": self.ports(options.ports)
                }
            ],
            "serve": {
                "name": name,
            }
        }
        if not os.path.exists('.metaparticle'):
            os.makedirs('.metaparticle')

        with open('.metaparticle/spec.json', 'w') as out:
            json.dump(svc, out)
    
        subprocess.check_call(['mp-compiler', '-f', '.metaparticle/spec.json'])