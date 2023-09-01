package common

//func Test_RemoveDuplicatesTagAssign_WhenNoDuplicatesExist(t *testing.T) {
//	// Arrange
//	input := []opslevel.TagInput{
//		{
//			Key:   "foo",
//			Value: "bar",
//		},
//		{
//			Key:   "hello",
//			Value: "world",
//		},
//		{
//			Key:   "apple",
//			Value: "orange",
//		},
//	}
//
//	// Act
//	result := removeDuplicatesFromTagInputList(input)
//
//	// Assert
//	autopilot.Equals(t, 3, len(result))
//}
//
//func Test_RemoveDuplicatesTagAssign_WhenDuplicatesExist(t *testing.T) {
//	// Arrange
//	input := []opslevel.TagInput{
//		{
//			Key:   "foo",
//			Value: "bar",
//		},
//		{
//			Key:   "hello",
//			Value: "world",
//		},
//		{
//			Key:   "foo",
//			Value: "bar",
//		},
//	}
//	// Act
//	result := removeDuplicatesFromTagInputList(input)
//
//	// Assert
//	autopilot.Equals(t, 2, len(result))
//}
