package io.metaparticle;

import io.metaparticle.annotations.Package;
import io.metaparticle.annotations.Runtime;
import java.io.IOException;
import java.lang.reflect.Method;
import java.nio.file.Files;
import java.nio.file.Paths;
import org.junit.After;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mockito;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import static io.metaparticle.Metaparticle.Containerize;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@RunWith(PowerMockRunner.class)
@PrepareForTest({Util.class})
public class MetaparticleTest {

	@After
	public void tearDown() throws IOException {
		Files.deleteIfExists(Paths.get("Dockerfile"));
	}

	@Test
	public void inDockerContainerTest() {
		assertEquals(false, Metaparticle.inDockerContainer());

	}

	@Runtime(iterations = 4, executor = "docker")
	@Package(repository = "test", verbose = true, jarFile = "test.jar")
	public static void dockerExecutorAnnotation() {
		Containerize(() -> {
			System.out.println("Hello batch job!");
		});
	}

	@Test
	public void getExecutorDockerTest() throws NoSuchMethodException, SecurityException {

		Method method = this.getClass().getMethod("dockerExecutorAnnotation");
		Runtime runtimeAnnotation = method.getAnnotation(Runtime.class);
		assertEquals(DockerImpl.class.getName(), Metaparticle.getExecutor(null).getClass().getName());
		assertEquals(DockerImpl.class.getName(), Metaparticle.getExecutor(runtimeAnnotation).getClass().getName());
	}

	@Runtime(iterations = 4, executor = "aci")
	public static void aciExecutorAnnotation() {
		Containerize(() -> {
			System.out.println("Hello batch job!");
		});
	}

	@Test
	public void getExecutorAciTest() throws NoSuchMethodException, SecurityException {
		Method method = this.getClass().getMethod("aciExecutorAnnotation");
		Runtime runtimeAnnotation = method.getAnnotation(Runtime.class);
		assertEquals(AciExecutor.class.getName(), Metaparticle.getExecutor(runtimeAnnotation).getClass().getName());

	}

	@Runtime(iterations = 4, executor = "metaparticle")
	public static void metaParticleExecutorAnnotation() {
		Containerize(() -> {
			System.out.println("Hello batch job!");
		});
	}

	@Test
	public void getExecutorMetaparticleTest() throws NoSuchMethodException, SecurityException {
		Method method = this.getClass().getMethod("metaParticleExecutorAnnotation");
		Runtime runtimeAnnotation = method.getAnnotation(Runtime.class);
		assertEquals(MetaparticleExecutor.class.getName(),
				Metaparticle.getExecutor(runtimeAnnotation).getClass().getName());

	}

	@Test
	public void getBuilderTest() throws NoSuchMethodException, SecurityException {
		Method method = this.getClass().getMethod("dockerExecutorAnnotation");
		Package packageAnnotation = method.getAnnotation(Package.class);
		assertEquals(DockerImpl.class.getName(), Metaparticle.getBuilder(null).getClass().getName());
		assertEquals(DockerImpl.class.getName(), Metaparticle.getBuilder(packageAnnotation).getClass().getName());
	}

	@Test
	public void dockerFileCommandSetForSpringBootApplication() throws IOException {
		Package mockPackage = mock(Package.class);
		when(mockPackage.dockerfile()).thenReturn("");
		when(mockPackage.jarFile()).thenReturn("target/foo.jar");
		Metaparticle.writeDockerfile("foo.Bar", mockPackage, true);

		String dockerfileAsString = new String(Files.readAllBytes(Paths.get("Dockerfile")));
		assertTrue(dockerfileAsString.endsWith("\nCMD java -jar /main.jar"));
	}

	@Test
	public void dockerFileCommandSetWithDefaultCommand() throws IOException {
		Package mockPackage = mock(Package.class);
		when(mockPackage.dockerfile()).thenReturn("");
		when(mockPackage.jarFile()).thenReturn("target/foo.jar");
		Metaparticle.writeDockerfile("foo.Bar", mockPackage, false);

		String dockerfileAsString = new String(Files.readAllBytes(Paths.get("Dockerfile")));
		assertTrue(dockerfileAsString.endsWith("\nCMD java -classpath /main.jar foo.Bar"));
	}

	@Test
	public void packageStepSkippedWhenRunningInSpringBootAppAndTargetAlreadyExists() {
		PowerMockito.mockStatic(Util.class);

		Metaparticle.doPackage(true, true);

		PowerMockito.verifyStatic(Util.class, Mockito.never());
		Util.handleErrorExec(new String[]{"mvn", "package"}, System.out, System.err);
	}

	@Test
	public void packageStepOccursWhenNotRunningInSpringBootAppAndTargetAlreadyExists() {
		PowerMockito.mockStatic(Util.class);

		Metaparticle.doPackage(false, true);

		PowerMockito.verifyStatic(Util.class);
		Util.handleErrorExec(new String[]{"mvn", "package"}, System.out, System.err);
	}

	@Test
	public void packageStepOccursWhenNotRunningInSpringBootAppAndTargetDoesNotExist() {
		PowerMockito.mockStatic(Util.class);

		Metaparticle.doPackage(false, false);

		PowerMockito.verifyStatic(Util.class);
		Util.handleErrorExec(new String[]{"mvn", "package"}, System.out, System.err);
	}

	@Test
	public void writeDockerFileTest() throws IOException, NoSuchMethodException, SecurityException {
		//TODO
	}

}
