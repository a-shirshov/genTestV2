package experiment

func twoSum(nums []int, target int) []int {
    hashMap := make(map[int]int)

    for index, number := range nums {
        val, ok := hashMap[number]
        if !ok {
            hashMap[target - number] = index
            continue
        }
        return []int{val,index}
    }

    return []int{}
}