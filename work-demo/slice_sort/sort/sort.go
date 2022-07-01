package sort

import "fmt"

// 冒泡排序--比较相邻的两个数大小，交换位置
func BubbleSort(arr []int) []int {
	for i := 0; i < len(arr); i++ {
		// 防止下标越位(len(arr)-i-1)
		for j := 0; j < len(arr)-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

// 选择排序--外循环从数组第0位开始，确定当前索引位置值为最小值，内循环不断与最小值比较，如果小于最小值则交换位置，一轮下来确定最左为最小值，然后外循环继续找到第二小的值
func SelectSort(arr []int) []int {
	for i := 0; i < len(arr); i++ {
		min := arr[i]
		for j := i; j < len(arr); j++ {
			if min > arr[j] {
				// 获取当前最小值
				min = arr[j]
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	return arr
}

// 插入排序--基本思想：将每一步的待排序的元素插入到一个已经排好序的序列中。可以理解为外循环就是序号的那个元素就是待排序元素，内循环做的事情就是将外循环的元素进行比较，插入到合适的位置
func InsertSort(arr []int) []int {
	for i := 0; i < len(arr); i++ {
		for j := i; j > 0; j-- {
			if arr[j] < arr[j-1] {
				arr[j-1], arr[j] = arr[j], arr[j-1]
			}
		}
	}
	return arr
}

//归并排序
func MergeSort(arr []int) []int {
	length := len(arr)
	if length < 2 {
		return arr
	}
	i := length / 2
	left := MergeSort(arr[0:i])
	right := MergeSort(arr[i:])
	res := merge(left, right)
	return res
}

//合并数组
func merge(left, right []int) []int {
	result := make([]int, 0)
	m, n := 0, 0
	l, r := len(left), len(right)
	//比较两个数组，谁小把元素值添加到结果集内
	for m < l && n < r {
		if left[m] > right[n] {
			result = append(result, right[n])
			n++
		} else {
			result = append(result, left[m])
			m++
		}
	}
	//如果有一个数组比完了，另一个数组还有元素的情况，则将剩余元素添加到结果集内
	result = append(result, right[n:]...)
	result = append(result, left[m:]...)
	return result
}

// 递归方法
func QuickSort(arr []int, p int, r int) []int {
	if p >= r {
		return arr
	}
	q := partition(arr, p, r)
	QuickSort(arr, p, q-1)
	QuickSort(arr, q+1, r)
	return arr
}

//排序并返回pivot
func partition(arr []int, p int, r int) int {
	k := arr[p]
	j := p
	for i := p; i < r; i++ {
		if k > arr[i] {
			arr[i], arr[j] = arr[j], arr[i]
			j++
		}
	}
	arr[r], arr[j] = arr[j], arr[r]
	return j
}

type student struct {
	Name string
	Age  int
}

// 地址传值
func ForSort() map[string]*student {
	silce := make(map[string]*student)
	array := make([]*student, 0)
	stus := []*student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}
	for _, stu := range stus {
		silce[stu.Name] = stu
		array = append(array, stu)
		fmt.Println("foreach：", silce[stu.Name])
	}
	for _, v := range array {
		fmt.Println("array: ", v)
	}
	return silce
}

//排序
func main() {
	arr := []int{5, 4, 12, 14, 21, 23, 9, 6, 8, 1, 3}
	list := BubbleSort(arr)
	fmt.Println("冒泡排序", list)
	listA := SelectSort(arr)
	fmt.Println("选择排序", listA)
	listB := InsertSort(arr)
	fmt.Println("插入排序", listB)
	listC := MergeSort(arr)
	fmt.Println("归并排序", listC)
	listD := QuickSort(listB, 2, 3)
	fmt.Println("递归方法", listD)
	listE := ForSort()
	fmt.Println("for循环地址", listE)
}
