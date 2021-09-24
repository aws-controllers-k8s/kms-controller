// At most one of Store ID and Store Name are allowed
if input.CustomKeyStoreId != nil && input.CustomKeyStoreName != nil {
  input.CustomKeyStoreName = nil
}