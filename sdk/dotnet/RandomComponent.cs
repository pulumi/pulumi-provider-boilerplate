// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Xyz
{
    [XyzResourceType("xyz:index:RandomComponent")]
    public partial class RandomComponent : global::Pulumi.ComponentResource
    {
        [Output("length")]
        public Output<int> Length { get; private set; } = null!;

        [Output("password")]
        public Output<string> Password { get; private set; } = null!;


        /// <summary>
        /// Create a RandomComponent resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public RandomComponent(string name, RandomComponentArgs args, ComponentResourceOptions? options = null)
            : base("xyz:index:RandomComponent", name, args ?? new RandomComponentArgs(), MakeResourceOptions(options, ""), remote: true)
        {
        }

        private static ComponentResourceOptions MakeResourceOptions(ComponentResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new ComponentResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = ComponentResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
    }

    public sealed class RandomComponentArgs : global::Pulumi.ResourceArgs
    {
        [Input("length", required: true)]
        public Input<int> Length { get; set; } = null!;

        public RandomComponentArgs()
        {
        }
        public static new RandomComponentArgs Empty => new RandomComponentArgs();
    }
}