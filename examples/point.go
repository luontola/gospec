// Copyright © 2009-2010 Esko Luontola <www.orfjackal.net>
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package examples


// Point2 implements the gospec.Equality interface, so it can be
// compared for equality in GoSpec specs.
//
// If the type can be used as a value, like in this case, then extra
// care is needed when writing the Equals method because the other
// object can be a value or a pointer to a value.
type Point2 struct {
	X, Y int
}

func (this Point2) Equals(other interface{}) bool {
	switch that := other.(type) {
	case Point2:
		return this.equals(&that)
	case *Point2:
		return this.equals(that)
	}
	return false
}

func (this *Point2) equals(that *Point2) bool {
	return this.X == that.X &&
	       this.Y == that.Y
}


// Point3 also implements the gospec.Equality interface, but unlike
// with Point2, here we assume that Point3 will not be used as a value.
// That makes the Equals method somewhat simpler.
type Point3 struct {
	X, Y, Z int
}

func (this *Point3) Equals(other interface{}) bool {
	switch that := other.(type) {
	case *Point3:
		return this.X == that.X &&
		       this.Y == that.Y &&
		       this.Z == that.Z
	}
	return false
}

