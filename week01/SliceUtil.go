package week01

// SliceDelIdx 删除slice中指定位置的元素
func SliceDelIdx[T any](slice *[]T, idx uint) {
	// 获取原slice的长度
	sliceLen := uint(len(*slice))
	// 判断是否越界, 或者slice的长度为0
	if idx >= sliceLen || sliceLen == 0 {
		return
	}
	switch idx {
	// 判断是否删除的是第一个元素 如果是则直接将slice的指针指向第二个元素
	case 0:
		*slice = (*slice)[1:]
	// 判断是否删除的是最后一个元素 如果是则直接将slice的长度减一
	case sliceLen - 1:
		*slice = (*slice)[:sliceLen-1]
	// 默认将idx后面的元素向前移动一位 并将slice的长度减一
	default:
		copy((*slice)[idx:], (*slice)[idx+1:])
		*slice = (*slice)[:sliceLen-1]
	}
	// 缩容
	ReduceCap(slice)
}

// ReduceCap 缩容slice
func ReduceCap[T any](slice *[]T) {
	// 扩容机制为：容量小于256时，扩容为原来的2倍 (cap = originCap * 2)，否则扩容为原来的1.25倍 (cap = originCap * 1.25)
	// 所以缩容机制可以设计为：容量大于256时，缩容为原来的0.75倍 (cap = originCap * 0.75)，否则0.5倍 (cap = originCap * 0.5)
	// 获取原slice的容量
	sliceCap := cap(*slice)
	sliceLen := len(*slice)
	threeQuartersCap := sliceCap * 3 / 4 // 0.75倍
	halfCap := sliceCap / 2              // 0.5倍
	var newSlice []T

	if sliceLen == 0 {
		return
	}
	// 判断是否需要缩容
	if sliceCap > 256 && sliceLen < threeQuartersCap { // 容量大于256且长度小于缩容后的容量
		newSlice = make([]T, sliceLen, threeQuartersCap)
	} else if sliceCap <= 256 && sliceLen < halfCap { // 容量小于256且长度小于缩容后的容量
		newSlice = make([]T, sliceLen, halfCap)
	} else {
		// 不需要缩容
		return
	}
	// 将原slice的元素复制到新slice中
	copy(newSlice, *slice)
	// 将新slice的指针赋值给原slice
	*slice = newSlice
}
