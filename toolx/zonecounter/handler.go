package zonecounter

import "sort"

type Zone struct {
	Begin float64
	End   float64
	Value float64
}

type Zones []Zone

func (z Zones) Add(zs Zones) Zones {
	z = append(z, zs...)
	return z
}

func (z Zones) DiscretizationWithMerge() Zones {
	if len(z) == 0 {
		return Zones{} // 返回已初始化的空切片
	}
	// 提取所有关键点
	var points []float64
	for _, zone := range z {
		points = append(points, zone.Begin, zone.End)
	}

	// 去重并排序
	points = uniqueAndSort(points)

	// 构造新的不相交区间
	var tempResult Zones
	for i := 0; i < len(points)-1; i++ {
		begin := points[i]
		end := points[i+1]

		// 计算当前区间的叠加值
		var value float64
		for _, zone := range z {
			if begin >= zone.Begin && end <= zone.End {
				value += zone.Value
			}
		}

		// 仅当叠加值大于 0 时，添加到临时结果
		if value > 0 {
			tempResult = append(tempResult, Zone{Begin: begin, End: end, Value: value})
		}
	}

	// 合并相邻且 Value 相同的区间
	var result Zones
	for i := 0; i < len(tempResult); i++ {
		if len(result) == 0 {
			result = append(result, tempResult[i])
		} else {
			last := &result[len(result)-1]
			if last.End == tempResult[i].Begin && last.Value == tempResult[i].Value {
				// 合并区间
				last.End = tempResult[i].End
			} else {
				result = append(result, tempResult[i])
			}
		}
	}

	return result
}

// 去重并排序
func uniqueAndSort(points []float64) []float64 {
	pointMap := make(map[float64]bool)
	for _, p := range points {
		pointMap[p] = true
	}

	var uniquePoints []float64
	for p := range pointMap {
		uniquePoints = append(uniquePoints, p)
	}

	sort.Float64s(uniquePoints)
	return uniquePoints
}
