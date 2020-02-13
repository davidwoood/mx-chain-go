package peer

// ListIndexUpdater will handle the updating of list type and the index for a peer
type ListIndexUpdater struct {
	updateListAndIndex func(pubKey string, list string, index int) error
}

// UpdateListAndIndex will update the list and the index for a given peer
func (liu *ListIndexUpdater) UpdateListAndIndex(pubKey string, list string, index int) error {
	return liu.updateListAndIndex(pubKey, list, index)
}

// IsInterfaceNil checks if the underlying object is nil
func (liu *ListIndexUpdater) IsInterfaceNil() bool {
	return liu == nil
}
