package loadinfo

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	SA_APS = "APS"
	SA_APC = "APC"
	SA_ASS = "ASS"
	SA_ASC = "ASC"
	SA_VPS = "VPS"
	SA_VPC = "VPC"
	SA_VSS = "VSS"
	SA_VSC = "VSC"

	BI_SRC = "SRC"
)

type StreamsAmt map[string]uint     //流数量
type LoadInfo map[string]StreamsAmt //负载信息

type StreamInfo struct {
	Si_APS uint `json:"streaming_aps"`
	Si_APC uint `json:"streaming_apc"`
	Si_ASS uint `json:"streaming_ass"`
	Si_ASC uint `json:"streaming_asc"`
	Si_VPS uint `json:"streaming_vps"`
	Si_VPC uint `json:"streaming_vpc"`
	Si_VSS uint `json:"streaming_vss"`
	Si_VSC uint `json:"streaming_vsc"`
}

type DetailInfo struct {
	Ipport     string     `json:"ipport"`
	Streaminfo StreamInfo `json:"streaminfo"`
}

type ResponseInfo struct {
	Total  StreamInfo   `json:"total"`
	Detail []DetailInfo `json:"detail"`
}

func (si *StreamInfo) Create(sa StreamsAmt) {
	si.Si_APS = sa[SA_APS]
	si.Si_APC = sa[SA_APC]
	si.Si_ASS = sa[SA_ASS]
	si.Si_ASC = sa[SA_ASC]
	si.Si_VPS = sa[SA_VPS]
	si.Si_VPC = sa[SA_VPC]
	si.Si_VSS = sa[SA_VSS]
	si.Si_VSC = sa[SA_VSC]
}

func (si *StreamInfo) Add(s *StreamInfo) {
	si.Si_VSS += s.Si_VSS
	si.Si_VSC += s.Si_VSC
	si.Si_VPS += s.Si_VPS
	si.Si_VPC += s.Si_VPC
	si.Si_ASC += s.Si_ASC
	si.Si_APC += s.Si_APC
	si.Si_APS += s.Si_APS
	si.Si_ASS += s.Si_ASS
}

func (si *StreamInfo) GetTotalSub() uint {
	return si.Si_VSS + si.Si_ASS + si.Si_ASC + si.Si_VSC
}

func (sa StreamsAmt) CheckKey(key string) bool {
	if key == SA_APS || key == SA_APC || key == SA_ASS || key == SA_ASC ||
		key == SA_VPS || key == SA_VPC || key == SA_VSS || key == SA_VSC {
		return true
	} else {
		return false
	}
}

func (sa StreamsAmt) Set(key string, val uint) bool {
	if !sa.CheckKey(key) {
		//键错误
		return false
	} else {
		sa[key] = val
		return true
	}
}

func (sa StreamsAmt) String() string {
	kwds := make([]string, 0)
	for k, v := range sa {
		tmp := fmt.Sprintf("%v:%v", k, v)
		kwds = append(kwds, tmp)
	}
	kwds_str := strings.Join(kwds, ", ")
	return fmt.Sprintf("{%v}", kwds_str)
}

func (sa StreamsAmt) FromString(s string) StreamsAmt {
	ret := make(StreamsAmt)
	ss := string([]rune(s)[1 : len(s)-1])
	kwds := strings.Split(ss, ", ")
	for _, kwd := range kwds {
		tmpKeyVal := strings.Split(kwd, ":")
		key, val := tmpKeyVal[0], tmpKeyVal[1]
		int_val, err := strconv.Atoi(val)
		if err != nil {
			//转换失败，数据有误
			return nil
		}
		if !ret.Set(key, uint(int_val)) {
			//赋值出错（键不存在）
			return nil
		}
	}
	return ret
}

// 解析请求中包含的负载信息
func (bi LoadInfo) Parse(msg string) bool {
	streamsAmount := make(StreamsAmt)
	words := strings.Split(msg, "&")
	for _, keyval := range words {
		tmpKeyVal := strings.Split(keyval, "=")
		key, val := tmpKeyVal[0], tmpKeyVal[1]
		if key == BI_SRC {
			bi[val] = streamsAmount
		} else {
			int_val, err := strconv.Atoi(val)
			if err != nil {
				//转换失败，数据有误
				return false
			}
			if !streamsAmount.Set(key, uint(int_val)) {
				//赋值出错（键不存在）
				return false
			}
		}
	}
	return true
}
