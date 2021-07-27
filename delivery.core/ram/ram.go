package ram

import (
	guuid "github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
)

type cMap map[string]models.Courier
type oMap map[string]models.Order
type tMap map[string]models.Task
type qMap map[string]models.Queue
type IDs string

type RAM struct {
	Couriers cMap
	Orders   oMap
	Tasks    tMap
	Queue    qMap
	IDs      IDs
}

func (ids *IDs) GenUUID() string {
	id := guuid.New()
	return id.String()
}

func (ids *IDs) SliceUUID(uuid string) string {
	return uuid[:5]
}

//func (qmap qMap) AddTasks(tasks []*models.Task) error {
//	var firstTask = *tasks[0]
//	var secondTask = *tasks[1]
//	var arrTasks models.Queue
//
//	for i := 0; i < len(qmap[firstTask.CourierUUID].Tasks); i++ {
//		arrTasks.Tasks = append(arrTasks.Tasks, qmap[firstTask.CourierUUID].Tasks[i])
//	}
//	if len(qmap[firstTask.CourierUUID].Tasks) != 0 {
//		lastIdx := len(qmap[firstTask.CourierUUID].Tasks) - 1
//		firstTask.ParentTaskUUID = qmap[firstTask.CourierUUID].Tasks[lastIdx].UUID
//	}
//	arrTasks.Tasks = append(arrTasks.Tasks, firstTask)
//	arrTasks.Tasks = append(arrTasks.Tasks, secondTask)
//
//	qmap[firstTask.CourierUUID] = arrTasks
//	return nil
//}

func (qmap qMap) FinishTask(courierUUID string) (error, string) {
	var queue models.Queue
	var finishedTaskUUID string
	queueLength := len(qmap[courierUUID].Tasks)

	if queueLength != 0 {
		if queueLength >= 2 {
			finishedTaskUUID = qmap[courierUUID].Tasks[0].UUID
			queue.Tasks = qmap[courierUUID].Tasks[1:]
		} else {
			finishedTaskUUID = qmap[courierUUID].Tasks[0].UUID
			queue.Tasks = []models.Task{}
		}
		qmap[courierUUID] = queue
		return nil, finishedTaskUUID
	}
	return errors.New("Current courier have no tasks to do"), finishedTaskUUID
}

//func (qmap qMap) DeleteTasks(courierUUID string, orderUUID string) error {
//	var newQueue models.Queue
//	var parentTaskUUID = "0"
//	var deletedTask int
//	for i := 0; i < len(qmap[courierUUID].Tasks); i++ {
//		if qmap[courierUUID].Tasks[i].OrderUUID != orderUUID {
//			newQueue.Tasks = append(newQueue.Tasks, qmap[courierUUID].Tasks[i])
//			newQueue.Tasks[i-deletedTask].ParentTaskUUID = parentTaskUUID
//			parentTaskUUID = newQueue.Tasks[i-deletedTask].UUID
//		} else {
//			deletedTask++
//		}
//	}
//	qmap[courierUUID] = newQueue
//	if deletedTask == 0 {
//		return errors.New("there is no tasks to delete")
//	}
//	return nil
//}

//func (qmap qMap) UpdateQueue(IDs *proto.QueueTasks) error {
//	queue := qmap[IDs.CourierUUID]
//	var task models.Task
//	var firstTaskIndex int
//	var secondTaskIndex int
//	for i := 0; i < len(qmap[IDs.CourierUUID].Tasks); i++ {
//		if qmap[IDs.CourierUUID].Tasks[i].UUID == IDs.FirstTaskUUID {
//			firstTaskIndex = i
//		}
//		if qmap[IDs.CourierUUID].Tasks[i].UUID == IDs.SecondTaskUUID {
//			secondTaskIndex = i
//		}
//	}
//	task = queue.Tasks[firstTaskIndex]
//	queue.Tasks[firstTaskIndex] = queue.Tasks[secondTaskIndex]
//	queue.Tasks[secondTaskIndex] = task
//	queue.Tasks[0].ParentTaskUUID = "0"
//	for i := 1; i < len(queue.Tasks); i++ {
//		queue.Tasks[i].ParentTaskUUID = queue.Tasks[i-1].UUID
//	}
//	qmap[IDs.CourierUUID] = queue
//	return nil
//}

func (qmap qMap) NextTask(courierUUID string) (error, models.Task) {
	var task models.Task
	lengthQueue := len(qmap[courierUUID].Tasks)
	if lengthQueue != 0 {
		task = qmap[courierUUID].Tasks[0]
		return nil, task
	}
	return errors.New("there is no tasks, try again later"), task
}
