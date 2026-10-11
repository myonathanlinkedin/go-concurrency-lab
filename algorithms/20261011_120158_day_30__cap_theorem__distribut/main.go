package main

func main() {
	// ---------- Test CP Mode ----------
	cp := NewCluster(ModeCP)
	if err := cp.Write("k", "v1"); err != nil {
		panic("CP write before partition should succeed")
	}
	if val, err := cp.Read("k"); err != nil || val != "v1" {
		panic("CP read before partition should return v1")
	}
	// Partition leaving only 1 alive node (majority lost)
	cp.Partition(2)
	if err := cp.Write("k", "v2"); err != ErrUnavailable {
		panic("CP write during minority partition should be unavailable")
	}
	if _, err := cp.Read("k"); err != ErrUnavailable {
		panic("CP read during minority partition should be unavailable")
	}
	cp.Heal()
	if err := cp.Write("k", "v2"); err != nil {
		panic("CP write after heal should succeed")
	}
	if val, err := cp.Read("k"); err != nil || val != "v2" {
		panic("CP read after heal should return v2")
	}

	// ---------- Test AP Mode ----------
	ap := NewCluster(ModeAP)
	if err := ap.Write("a", "1"); err != nil {
		panic("AP write before partition should succeed")
	}
	if val, err := ap.Read("a"); err != nil || val != "1" {
		panic("AP read before partition should return 1")
	}
	// Partition leaving only 1 alive node
	ap.Partition(2)
	if err := ap.Write("a", "2"); err != nil {
		panic("AP write during partition should succeed despite minority")
	}
	if val, err := ap.Read("a"); err != nil || val != "2" {
		panic("AP read during partition should return the latest local value")
	}
	ap.Heal()
	// After heal, the healed nodes may still hold old value; AP does not guarantee consistency.
	if val, err := ap.Read("a"); err != nil {
		panic("AP read after heal should succeed")
	} else {
		// Value may be either "1" or "2" depending on which node responded.
		if val != "1" && val != "2" {
			panic("AP read after heal returned unexpected value")
		}
	}

	// ---------- Test CA Mode ----------
	ca := NewCluster(ModeCA)
	if err := ca.Write("x", "alpha"); err != nil {
		panic("CA write before partition should succeed")
	}
	if val, err := ca.Read("x"); err != nil || val != "alpha" {
		panic("CA read before partition should return alpha")
	}
	// Attempt to partition; CA should reject operations.
	ca.Partition(1)
	if err := ca.Write("x", "beta"); err != ErrInvalidMode {
		panic("CA write during partition should be invalid")
	}
	if _, err := ca.Read("x"); err != ErrInvalidMode {
		panic("CA read during partition should be invalid")
	}
	// Heal and verify normal operation resumes.
	ca.Heal()
	if err := ca.Write("x", "beta"); err != nil {
		panic("CA write after heal should succeed")
	}
	if val, err := ca.Read("x"); err != nil || val != "beta" {
		panic("CA read after heal should return beta")
	}
}
