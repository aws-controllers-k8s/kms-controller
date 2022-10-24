    policy, err := rm.getPolicy(ctx, &resource{ko})
    if err != nil {
        return &resource{ko}, err
    }
    ko.Spec.Policy = policy
    tags, err := rm.listTags(ctx, &resource{ko})
    if err != nil {
        return &resource{ko}, err
    }
    ko.Spec.Tags = FromACKTags(tags)
    keyRotationStatus, err := rm.getKeyRotationStatus(&resource{ko})
	if err != nil {
		return &resource{ko}, err
	}
	ko.Spec.EnableKeyRotation = keyRotationStatus.KeyRotationEnabled