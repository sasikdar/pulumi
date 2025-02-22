// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Example
{
    [ExampleResourceType("example::Component")]
    public partial class Component : Pulumi.CustomResource
    {
        [Output("provider")]
        public Output<Pulumi.Kubernetes.Provider?> Provider { get; private set; } = null!;

        [Output("securityGroup")]
        public Output<Pulumi.Aws.Ec2.SecurityGroup?> SecurityGroup { get; private set; } = null!;

        [Output("storageClasses")]
        public Output<ImmutableDictionary<string, Pulumi.Kubernetes.Storage.V1.StorageClass>?> StorageClasses { get; private set; } = null!;


        /// <summary>
        /// Create a Component resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public Component(string name, ComponentArgs? args = null, CustomResourceOptions? options = null)
            : base("example::Component", name, args ?? new ComponentArgs(), MakeResourceOptions(options, ""))
        {
        }

        private Component(string name, Input<string> id, CustomResourceOptions? options = null)
            : base("example::Component", name, null, MakeResourceOptions(options, id))
        {
        }

        private static CustomResourceOptions MakeResourceOptions(CustomResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new CustomResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = CustomResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
        /// <summary>
        /// Get an existing Component resource's state with the given name, ID, and optional extra
        /// properties used to qualify the lookup.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resulting resource.</param>
        /// <param name="id">The unique provider ID of the resource to lookup.</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public static Component Get(string name, Input<string> id, CustomResourceOptions? options = null)
        {
            return new Component(name, id, options);
        }
    }

    public sealed class ComponentArgs : Pulumi.ResourceArgs
    {
        [Input("metadata")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs>? Metadata { get; set; }

        public ComponentArgs()
        {
        }
    }
}
