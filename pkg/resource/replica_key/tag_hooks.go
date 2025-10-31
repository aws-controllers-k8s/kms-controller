// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package replica_key

import (
	"context"

	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/kms"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
)

// updateTags performs the TagResource API call using Spec.Tags field of
// resource in the parameter
func (rm *resourceManager) updateTags(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.updateTags")
	defer func() {
		exit(err)
	}()
	latestTags, err := rm.listTags(ctx, r)
	if err != nil {
		return err
	}
	desiredTags, _ := convertToOrderedACKTags(r.ko.Spec.Tags)
	// First remove the keys that are not present in desired state anymore
	tagKeysToRemove := removedTagKeys(desiredTags, latestTags)
	if tagKeysToRemove != nil && len(tagKeysToRemove) > 0 {
		untagKeyInput := svcsdk.UntagResourceInput{
			KeyId:   r.ko.Status.ReplicaKeyMetadata.KeyID,
			TagKeys: derefStringSlice(tagKeysToRemove),
		}
		_, err = rm.sdkapi.UntagResource(ctx, &untagKeyInput)
		rm.metrics.RecordAPICall("UPDATE", "UntagResource", err)
		if err != nil {
			return err
		}
	}
	// Now tag the KMS Key with desired tags
	if len(desiredTags) == 0 {
		return nil
	}
	var svcTags []svcsdktypes.Tag
	for k, v := range desiredTags {
		kCopy := k
		vCopy := v
		tag := svcsdktypes.Tag{
			TagKey:   &kCopy,
			TagValue: &vCopy,
		}
		svcTags = append(svcTags, tag)
	}
	tagKeyInput := svcsdk.TagResourceInput{
		KeyId: r.ko.Status.ReplicaKeyMetadata.KeyID,
		Tags:  svcTags,
	}
	_, err = rm.sdkapi.TagResource(ctx, &tagKeyInput)
	rm.metrics.RecordAPICall("UPDATE", "TagResource", err)
	return err
}

// listTags performs the ListResourceTags API call and returns the result in
// form of acktags.Tags format
func (rm *resourceManager) listTags(ctx context.Context, r *resource) (tags acktags.Tags, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.listTags")
	defer func() {
		exit(err)
	}()
	var truncated = true
	var marker *string
	tags = acktags.NewTags()
	for truncated {
		listTagsInput := svcsdk.ListResourceTagsInput{
			KeyId:  r.ko.Status.ReplicaKeyMetadata.KeyID,
			Marker: marker,
		}
		resp, err := rm.sdkapi.ListResourceTags(ctx, &listTagsInput)
		rm.metrics.RecordAPICall("GET", "ListResourceTags", err)
		if err != nil {
			return nil, err
		}
		if !resp.Truncated {
			truncated = false
		} else {
			truncated = resp.Truncated
		}
		marker = resp.NextMarker
		for _, t := range resp.Tags {
			tags[*t.TagKey] = *t.TagValue
		}
	}
	return tags, nil
}

// removedTagKeys returns the tag keys that are present inside latestTags but
// are not part of desiredTags
func removedTagKeys(desiredTags acktags.Tags, latestTags acktags.Tags) []*string {
	var removedKeys []*string
	for k := range latestTags {
		if _, found := desiredTags[k]; !found {
			kCopy := k
			removedKeys = append(removedKeys, &kCopy)
		}
	}
	return removedKeys
}

func derefStringSlice(in []*string) []string {
	out := make([]string, len(in))
	for i, s := range in {
		if s != nil {
			out[i] = *s
		}
	}
	return out
}
