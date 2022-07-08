	// TODO(vijtrip2): remove this pagination handling once it is handled by the
	// ACK code-generator. https://github.com/aws-controllers-k8s/community/issues/1383
	aliases := []*svcsdk.AliasListEntry{}
	aliases = append(aliases, resp.Aliases...)
	for resp.Truncated != nil && *resp.Truncated {
		input.Marker = resp.NextMarker
		resp, err = rm.sdkapi.ListAliasesWithContext(ctx, input)
		rm.metrics.RecordAPICall("READ_MANY", "ListAliases", err)
		if err != nil {
			if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "UNKNOWN" {
				return nil, ackerr.NotFound
			}
			return nil, err
		}
		aliases = append(aliases, resp.Aliases...)
	}
	// Filter resulting aliases, matching only the one with the name in the spec
	matchingAliases := []*svcsdk.AliasListEntry{}
	for _, elem := range aliases {
	  if elem.AliasName == nil || r.ko.Spec.Name == nil {
		continue
	  }

	  if *elem.AliasName == *r.ko.Spec.Name {
		matchingAliases = append(matchingAliases, elem)
	  }
	}
	resp.Aliases = matchingAliases