const mp = require('../index');
const sinon = require('sinon');
const dockerBuilder = require('../docker-builder');
const dockerRunner = require('../docker-runner');
const fs = require('fs');

describe('containerize', () => {
    let buildStub, runStub, readFileSyncStub, writeFileSyncStub;

    beforeEach(() => {
        buildStub = sinon.stub(dockerBuilder, 'build');
        runStub = sinon.stub(dockerRunner, 'run');
        readFileSyncStub = sinon.stub(fs, 'readFileSync');
        writeFileSyncStub = sinon.stub(fs, 'writeFileSync');
    })
    afterEach(() => {
        process.env.METAPARTICLE_IN_CONTAINER = undefined;
        buildStub.restore();
        runStub.restore();
        readFileSyncStub.restore();
        writeFileSyncStub.restore();
    })

    describe('METAPARTICLE_IN_CONTAINER', () => {
        it('should execute function when METAPARTICLE_IN_CONTAINER is "1"', () => {
            process.env.METAPARTICLE_IN_CONTAINER = '1';
            let func = sinon.spy();
            mp.containerize(func);
            sinon.assert.called(func);
        })
        it('should execute function when METAPARTICLE_IN_CONTAINER is "true"', () => {
            process.env.METAPARTICLE_IN_CONTAINER = 'true';
            let func = sinon.spy();
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

    describe('options', () => {
        let publishStub;
        let pkgName = 'pkgName';

        beforeEach(() => {
            readFileSyncStub.callsFake(() => `{"name":"${pkgName}"}`);
            publishStub = sinon.stub(dockerBuilder, 'publish');
            process.env.METAPARTICLE_IN_CONTAINER = 'false';
        })
        afterEach(() => {
            process.env.METAPARTICLE_IN_CONTAINER = undefined;            
            publishStub.restore();
        })

        it('should publish image', () => {
            mp.containerize({ publish: true });
            sinon.assert.calledWith(publishStub, pkgName);
        })
        it('should publish image in given repository', () => {
            let repository = 'repository';
            mp.containerize({ publish: true, repository });
            sinon.assert.calledWith(publishStub, repository + '/' + pkgName);
        })
        it('should write dockerfile with given name', () => {
            let name = 'name';
            mp.containerize({ name });
            sinon.assert.calledWithMatch(writeFileSyncStub, 'Dockerfile', `COPY ./ /${name}/`)
        })
    })
})