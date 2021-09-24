// Filter resulting aliases, matching only the one with the name in the spec
matchingAliases := []*svcsdk.AliasListEntry{}
for _, elem := range resp.Aliases {
  if *elem.AliasName == *r.ko.Spec.Name {
    matchingAliases = append(matchingAliases, elem)
  }
}
resp.Aliases = matchingAliases