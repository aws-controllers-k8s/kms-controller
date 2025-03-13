    policy, err := rm.getPolicy(ctx, &resource{ko})
    if err != nil {
        return &resource{ko}, err
    }
    ko.Spec.Policy = policy
    tags, err := rm.listTags(ctx, &resource{ko})
    if err != nil {
        return &resource{ko}, err
    }
    ko.Spec.Tags = fromACKTags(tags, nil)
    keyRotationStatus, err := rm.getKeyRotationStatus(ctx, &resource{ko})
	if err != nil || keyRotationStatus == nil {
		return &resource{ko}, err
	}
	enabled := keyRotationStatus.KeyRotationEnabled
	ko.Spec.EnableKeyRotation = &enabled