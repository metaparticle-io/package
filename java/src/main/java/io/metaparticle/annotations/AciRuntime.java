package io.metaparticle.annotations;

import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;

@Retention(RetentionPolicy.RUNTIME)
public @interface AciRuntime {
    /**
     * Azure tenant id.
     */
    public String azureTenantId() default "";

    /**
     * Azure subscription id.
     * Container instances will be created using the given subscription.
     */
    public String azureSubscriptionId() default "";

    /**
     * Azure client id (aka. application id).
     */
    public String azureClientId() default "";

    /**
     * Azure client secret (aka. application key).
     * If specified, then will authenticate to Azure as service principal.
     * And then no need to specify username/password.
     */
    public String azureClientSecret() default "";

    public String username() default "";

    public String password() default "";

    /**
     * The resource group to create container instance.
     */
    public String aciResourceGroup() default "metaparticle-execution";
}
