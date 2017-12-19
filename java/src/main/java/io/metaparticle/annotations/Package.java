package io.metaparticle.annotations;

import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;

@Retention(RetentionPolicy.RUNTIME)
public @interface Package {
    String repository() default "";

    boolean verbose() default false;

    boolean quiet() default false;

    String version() default "";

    String builder() default "docker";

    boolean build() default true;

    boolean publish() default false;

    String jarFile() default "";

    String dockerfile() default "";
}
