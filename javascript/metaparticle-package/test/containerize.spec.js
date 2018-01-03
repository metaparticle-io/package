var mp = require('../index');
var sinon = require('sinon');
var dockerBuilder = require('../docker-builder');
var dockerRunner = require('../docker-runner');
var fs = require('fs');

describe('containerize', () => {
    describe('METAPARTICLE_IN_CONTAINER', () => {
        var buildStub, runStub, readFileSyncStub;
        beforeEach(() => {
            buildStub = sinon.stub(dockerBuilder, 'build');
            runStub = sinon.stub(dockerRunner, 'run');
            readFileSyncStub = sinon.stub(fs, 'readFileSync');
            sinon.stub(fs, 'writeFileSync');
        })
        afterEach(() => {
            process.env.METAPARTICLE_IN_CONTAINER = undefined;
            buildStub.restore();
            runStub.restore();
            readFileSyncStub.restore();
            fs.writeFileSync.restore();
        })
        it('should execute function when METAPARTICLE_IN_CONTAINER is "1"', () => {
            process.env.METAPARTICLE_IN_CONTAINER = '1';
            var func = sinon.spy();
            mp.containerize(func);
            sinon.assert.called(func);
        })
        it('should execute function when METAPARTICLE_IN_CONTAINER is "true"', () => {
            process.env.METAPARTICLE_IN_CONTAINER = 'true';
            var func = sinon.spy();
            mp.containerize(func);
            sinon.assert.called(func);
        })
        it('should build image and run container when METAPARTICLE_IN_CONTAINER is "0"', () => {
            process.env.METAPARTICLE_IN_CONTAINER = '0';
            readFileSyncStub.callsFake(() => "{}");
            mp.containerize({});
            sinon.assert.called(buildStub);
            sinon.assert.called(runStub);
        })
        it('should build image and run container when METAPARTICLE_IN_CONTAINER is "false"', () => {
            process.env.METAPARTICLE_IN_CONTAINER = 'false';
            readFileSyncStub.callsFake(() => "{}");
            mp.containerize({});
            sinon.assert.called(buildStub);
            sinon.assert.called(runStub);
        })
    })
})