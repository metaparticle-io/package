package io.metaparticle.test;

import static io.metaparticle.Metaparticle.Containerize;
import static org.junit.jupiter.api.Assertions.assertEquals;

import java.io.IOException;
import java.lang.reflect.Method;
import java.nio.file.Files;
import java.nio.file.Paths;

import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.powermock.api.mockito.PowerMockito;

import io.metaparticle.AciExecutor;
import io.metaparticle.DockerImpl;
import io.metaparticle.Metaparticle;
import io.metaparticle.MetaparticleExecutor;
import io.metaparticle.annotations.Package;
import io.metaparticle.annotations.Runtime;

public class MetaparticleTest {

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
	public void writeDockerFileTest() throws IOException, NoSuchMethodException, SecurityException {
		//TODO
	}

}
