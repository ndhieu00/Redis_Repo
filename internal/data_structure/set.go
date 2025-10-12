package data_structure

// Set represents a Redis set (unique string collection)
type Set map[string]struct{}

// NewSet creates a new set with initial members
func NewSet(members []string) Set {
	set := make(Set)
	for _, member := range members {
		set[member] = struct{}{}
	}
	return set
}

// Add adds members to the set, returns number of new members added
func (s Set) Add(members []string) int {
	if s == nil {
		return 0
	}

	added := 0
	for _, member := range members {
		if _, exists := s[member]; !exists {
			s[member] = struct{}{}
			added++
		}
	}
	return added
}

// Remove removes members from the set, returns number of members removed
func (s Set) Remove(members []string) int {
	if s == nil {
		return 0
	}

	removed := 0
	for _, member := range members {
		if _, exists := s[member]; exists {
			delete(s, member)
			removed++
		}
	}

	return removed
}

// IsMember checks if a member exists in the set
func (s Set) IsMember(member string) int {
	if s == nil {
		return 0
	}

	_, exist := s[member]
	if exist {
		return 1
	}
	return 0
}
