const mp = require('@metaparticle/package');

mp.containerize({
        jobSpec: {
            count: 4
        },
        runner: 'metaparticle',
        repository: 'brendanburns',
        publish: true,
    },
    () => {
        console.log("I am a batch job!");
    }
);