package io.metaparticle.annotations;

import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;

@Retention(RetentionPolicy.RUNTIME)
public @interface Runtime {
    public int replicas() default 1;

    public int shards() default 0;

    public String urlShardPattern() default "";

    public String executor() default "docker";

    public int[] ports() default {};

    public boolean publicAddress() default false;

    public boolean election() default false;
}
