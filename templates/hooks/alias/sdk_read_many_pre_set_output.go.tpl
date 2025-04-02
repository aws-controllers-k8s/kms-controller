	// TODO(vijtrip2): remove this pagination handling once it is handled by the
	// ACK code-generator. https://github.com/aws-controllers-k8s/community/issues/1383
	aliases := []svcsdktypes.AliasListEntry{}
	aliases = append(aliases, resp.Aliases...)
	for resp.Truncated {
		input.Marker = resp.NextMarker
		resp, err = rm.sdkapi.ListAliases(ctx, input)
		rm.metrics.RecordAPICall("READ_MANY", "ListAliases", err)
		if err != nil {
			if awsErr, ok := ackerr.AWSError(err); ok && awsErr.ErrorCode() == "NotFound" {
				return nil, ackerr.NotFound
			}
			return nil, err
		}
		aliases = append(aliases, resp.Aliases...)
	}
	aliasName := ensureAliasName(r.ko.Spec.Name)
	// Filter resulting aliases, matching only the one with the name in the spec
	matchingAliases := []svcsdktypes.AliasListEntry{}
	for _, elem := range aliases {
		if elem.AliasName == nil || aliasName == nil {
			continue
		}

		if *elem.AliasName == *aliasName {
			matchingAliases = append(matchingAliases, elem)
		}
	}
	resp.Aliases = matchingAliases