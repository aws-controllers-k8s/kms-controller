    policy, err := rm.getPolicy(ctx, &resource{ko})
    if err != nil {
        return &resource{ko}, err
    }
    ko.Spec.Policy = policy