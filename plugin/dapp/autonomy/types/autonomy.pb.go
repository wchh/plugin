// Code generated by protoc-gen-go. DO NOT EDIT.
// source: autonomy.proto

package types

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// message for execs.Autonomy
type AutonomyAction struct {
	// Types that are valid to be assigned to Value:
	//	*AutonomyAction_PropBoard
	//	*AutonomyAction_RvkPropBoard
	//	*AutonomyAction_VotePropBoard
	//	*AutonomyAction_TmintPropBoard
	//	*AutonomyAction_PropProject
	//	*AutonomyAction_RvkPropProject
	//	*AutonomyAction_VotePropProject
	//	*AutonomyAction_PubVotePropProject
	//	*AutonomyAction_TmintPropProject
	//	*AutonomyAction_PropRule
	//	*AutonomyAction_RvkPropRule
	//	*AutonomyAction_VotePropRule
	//	*AutonomyAction_TmintPropRule
	//	*AutonomyAction_Transfer
	//	*AutonomyAction_CommentProp
	//	*AutonomyAction_PropChange
	//	*AutonomyAction_RvkPropChange
	//	*AutonomyAction_VotePropChange
	//	*AutonomyAction_TmintPropChange
	Value                isAutonomyAction_Value `protobuf_oneof:"value"`
	Ty                   int32                  `protobuf:"varint,20,opt,name=ty,proto3" json:"ty,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *AutonomyAction) Reset()         { *m = AutonomyAction{} }
func (m *AutonomyAction) String() string { return proto.CompactTextString(m) }
func (*AutonomyAction) ProtoMessage()    {}
func (*AutonomyAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_0246b47df8434d60, []int{0}
}

func (m *AutonomyAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AutonomyAction.Unmarshal(m, b)
}
func (m *AutonomyAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AutonomyAction.Marshal(b, m, deterministic)
}
func (m *AutonomyAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AutonomyAction.Merge(m, src)
}
func (m *AutonomyAction) XXX_Size() int {
	return xxx_messageInfo_AutonomyAction.Size(m)
}
func (m *AutonomyAction) XXX_DiscardUnknown() {
	xxx_messageInfo_AutonomyAction.DiscardUnknown(m)
}

var xxx_messageInfo_AutonomyAction proto.InternalMessageInfo

type isAutonomyAction_Value interface {
	isAutonomyAction_Value()
}

type AutonomyAction_PropBoard struct {
	PropBoard *ProposalBoard `protobuf:"bytes,1,opt,name=propBoard,proto3,oneof"`
}

type AutonomyAction_RvkPropBoard struct {
	RvkPropBoard *RevokeProposalBoard `protobuf:"bytes,2,opt,name=rvkPropBoard,proto3,oneof"`
}

type AutonomyAction_VotePropBoard struct {
	VotePropBoard *VoteProposalBoard `protobuf:"bytes,3,opt,name=votePropBoard,proto3,oneof"`
}

type AutonomyAction_TmintPropBoard struct {
	TmintPropBoard *TerminateProposalBoard `protobuf:"bytes,4,opt,name=tmintPropBoard,proto3,oneof"`
}

type AutonomyAction_PropProject struct {
	PropProject *ProposalProject `protobuf:"bytes,5,opt,name=propProject,proto3,oneof"`
}

type AutonomyAction_RvkPropProject struct {
	RvkPropProject *RevokeProposalProject `protobuf:"bytes,6,opt,name=rvkPropProject,proto3,oneof"`
}

type AutonomyAction_VotePropProject struct {
	VotePropProject *VoteProposalProject `protobuf:"bytes,7,opt,name=votePropProject,proto3,oneof"`
}

type AutonomyAction_PubVotePropProject struct {
	PubVotePropProject *PubVoteProposalProject `protobuf:"bytes,8,opt,name=pubVotePropProject,proto3,oneof"`
}

type AutonomyAction_TmintPropProject struct {
	TmintPropProject *TerminateProposalProject `protobuf:"bytes,9,opt,name=tmintPropProject,proto3,oneof"`
}

type AutonomyAction_PropRule struct {
	PropRule *ProposalRule `protobuf:"bytes,10,opt,name=propRule,proto3,oneof"`
}

type AutonomyAction_RvkPropRule struct {
	RvkPropRule *RevokeProposalRule `protobuf:"bytes,11,opt,name=rvkPropRule,proto3,oneof"`
}

type AutonomyAction_VotePropRule struct {
	VotePropRule *VoteProposalRule `protobuf:"bytes,12,opt,name=votePropRule,proto3,oneof"`
}

type AutonomyAction_TmintPropRule struct {
	TmintPropRule *TerminateProposalRule `protobuf:"bytes,13,opt,name=tmintPropRule,proto3,oneof"`
}

type AutonomyAction_Transfer struct {
	Transfer *TransferFund `protobuf:"bytes,14,opt,name=transfer,proto3,oneof"`
}

type AutonomyAction_CommentProp struct {
	CommentProp *Comment `protobuf:"bytes,15,opt,name=commentProp,proto3,oneof"`
}

type AutonomyAction_PropChange struct {
	PropChange *ProposalChange `protobuf:"bytes,16,opt,name=propChange,proto3,oneof"`
}

type AutonomyAction_RvkPropChange struct {
	RvkPropChange *RevokeProposalChange `protobuf:"bytes,17,opt,name=rvkPropChange,proto3,oneof"`
}

type AutonomyAction_VotePropChange struct {
	VotePropChange *VoteProposalChange `protobuf:"bytes,18,opt,name=votePropChange,proto3,oneof"`
}

type AutonomyAction_TmintPropChange struct {
	TmintPropChange *TerminateProposalChange `protobuf:"bytes,19,opt,name=tmintPropChange,proto3,oneof"`
}

func (*AutonomyAction_PropBoard) isAutonomyAction_Value() {}

func (*AutonomyAction_RvkPropBoard) isAutonomyAction_Value() {}

func (*AutonomyAction_VotePropBoard) isAutonomyAction_Value() {}

func (*AutonomyAction_TmintPropBoard) isAutonomyAction_Value() {}

func (*AutonomyAction_PropProject) isAutonomyAction_Value() {}

func (*AutonomyAction_RvkPropProject) isAutonomyAction_Value() {}

func (*AutonomyAction_VotePropProject) isAutonomyAction_Value() {}

func (*AutonomyAction_PubVotePropProject) isAutonomyAction_Value() {}

func (*AutonomyAction_TmintPropProject) isAutonomyAction_Value() {}

func (*AutonomyAction_PropRule) isAutonomyAction_Value() {}

func (*AutonomyAction_RvkPropRule) isAutonomyAction_Value() {}

func (*AutonomyAction_VotePropRule) isAutonomyAction_Value() {}

func (*AutonomyAction_TmintPropRule) isAutonomyAction_Value() {}

func (*AutonomyAction_Transfer) isAutonomyAction_Value() {}

func (*AutonomyAction_CommentProp) isAutonomyAction_Value() {}

func (*AutonomyAction_PropChange) isAutonomyAction_Value() {}

func (*AutonomyAction_RvkPropChange) isAutonomyAction_Value() {}

func (*AutonomyAction_VotePropChange) isAutonomyAction_Value() {}

func (*AutonomyAction_TmintPropChange) isAutonomyAction_Value() {}

func (m *AutonomyAction) GetValue() isAutonomyAction_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *AutonomyAction) GetPropBoard() *ProposalBoard {
	if x, ok := m.GetValue().(*AutonomyAction_PropBoard); ok {
		return x.PropBoard
	}
	return nil
}

func (m *AutonomyAction) GetRvkPropBoard() *RevokeProposalBoard {
	if x, ok := m.GetValue().(*AutonomyAction_RvkPropBoard); ok {
		return x.RvkPropBoard
	}
	return nil
}

func (m *AutonomyAction) GetVotePropBoard() *VoteProposalBoard {
	if x, ok := m.GetValue().(*AutonomyAction_VotePropBoard); ok {
		return x.VotePropBoard
	}
	return nil
}

func (m *AutonomyAction) GetTmintPropBoard() *TerminateProposalBoard {
	if x, ok := m.GetValue().(*AutonomyAction_TmintPropBoard); ok {
		return x.TmintPropBoard
	}
	return nil
}

func (m *AutonomyAction) GetPropProject() *ProposalProject {
	if x, ok := m.GetValue().(*AutonomyAction_PropProject); ok {
		return x.PropProject
	}
	return nil
}

func (m *AutonomyAction) GetRvkPropProject() *RevokeProposalProject {
	if x, ok := m.GetValue().(*AutonomyAction_RvkPropProject); ok {
		return x.RvkPropProject
	}
	return nil
}

func (m *AutonomyAction) GetVotePropProject() *VoteProposalProject {
	if x, ok := m.GetValue().(*AutonomyAction_VotePropProject); ok {
		return x.VotePropProject
	}
	return nil
}

func (m *AutonomyAction) GetPubVotePropProject() *PubVoteProposalProject {
	if x, ok := m.GetValue().(*AutonomyAction_PubVotePropProject); ok {
		return x.PubVotePropProject
	}
	return nil
}

func (m *AutonomyAction) GetTmintPropProject() *TerminateProposalProject {
	if x, ok := m.GetValue().(*AutonomyAction_TmintPropProject); ok {
		return x.TmintPropProject
	}
	return nil
}

func (m *AutonomyAction) GetPropRule() *ProposalRule {
	if x, ok := m.GetValue().(*AutonomyAction_PropRule); ok {
		return x.PropRule
	}
	return nil
}

func (m *AutonomyAction) GetRvkPropRule() *RevokeProposalRule {
	if x, ok := m.GetValue().(*AutonomyAction_RvkPropRule); ok {
		return x.RvkPropRule
	}
	return nil
}

func (m *AutonomyAction) GetVotePropRule() *VoteProposalRule {
	if x, ok := m.GetValue().(*AutonomyAction_VotePropRule); ok {
		return x.VotePropRule
	}
	return nil
}

func (m *AutonomyAction) GetTmintPropRule() *TerminateProposalRule {
	if x, ok := m.GetValue().(*AutonomyAction_TmintPropRule); ok {
		return x.TmintPropRule
	}
	return nil
}

func (m *AutonomyAction) GetTransfer() *TransferFund {
	if x, ok := m.GetValue().(*AutonomyAction_Transfer); ok {
		return x.Transfer
	}
	return nil
}

func (m *AutonomyAction) GetCommentProp() *Comment {
	if x, ok := m.GetValue().(*AutonomyAction_CommentProp); ok {
		return x.CommentProp
	}
	return nil
}

func (m *AutonomyAction) GetPropChange() *ProposalChange {
	if x, ok := m.GetValue().(*AutonomyAction_PropChange); ok {
		return x.PropChange
	}
	return nil
}

func (m *AutonomyAction) GetRvkPropChange() *RevokeProposalChange {
	if x, ok := m.GetValue().(*AutonomyAction_RvkPropChange); ok {
		return x.RvkPropChange
	}
	return nil
}

func (m *AutonomyAction) GetVotePropChange() *VoteProposalChange {
	if x, ok := m.GetValue().(*AutonomyAction_VotePropChange); ok {
		return x.VotePropChange
	}
	return nil
}

func (m *AutonomyAction) GetTmintPropChange() *TerminateProposalChange {
	if x, ok := m.GetValue().(*AutonomyAction_TmintPropChange); ok {
		return x.TmintPropChange
	}
	return nil
}

func (m *AutonomyAction) GetTy() int32 {
	if m != nil {
		return m.Ty
	}
	return 0
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*AutonomyAction) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*AutonomyAction_PropBoard)(nil),
		(*AutonomyAction_RvkPropBoard)(nil),
		(*AutonomyAction_VotePropBoard)(nil),
		(*AutonomyAction_TmintPropBoard)(nil),
		(*AutonomyAction_PropProject)(nil),
		(*AutonomyAction_RvkPropProject)(nil),
		(*AutonomyAction_VotePropProject)(nil),
		(*AutonomyAction_PubVotePropProject)(nil),
		(*AutonomyAction_TmintPropProject)(nil),
		(*AutonomyAction_PropRule)(nil),
		(*AutonomyAction_RvkPropRule)(nil),
		(*AutonomyAction_VotePropRule)(nil),
		(*AutonomyAction_TmintPropRule)(nil),
		(*AutonomyAction_Transfer)(nil),
		(*AutonomyAction_CommentProp)(nil),
		(*AutonomyAction_PropChange)(nil),
		(*AutonomyAction_RvkPropChange)(nil),
		(*AutonomyAction_VotePropChange)(nil),
		(*AutonomyAction_TmintPropChange)(nil),
	}
}

func init() {
	proto.RegisterType((*AutonomyAction)(nil), "types.AutonomyAction")
}

func init() {
	proto.RegisterFile("autonomy.proto", fileDescriptor_0246b47df8434d60)
}

var fileDescriptor_0246b47df8434d60 = []byte{
	// 504 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x94, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x40, 0xd3, 0x40, 0xfa, 0x31, 0x8e, 0x9d, 0x32, 0x2d, 0x60, 0xc2, 0x57, 0xc5, 0xa9, 0xa7,
	0x48, 0x14, 0x24, 0x24, 0xa4, 0x4a, 0x6d, 0x83, 0x42, 0x84, 0x84, 0x88, 0xac, 0xaa, 0x77, 0x27,
	0x2c, 0x10, 0x1a, 0x7b, 0xad, 0xcd, 0x3a, 0x52, 0x7e, 0x00, 0xff, 0x1b, 0x79, 0x3c, 0x5e, 0x7b,
	0x37, 0xee, 0x2d, 0xf6, 0xce, 0x7b, 0xd1, 0x3c, 0x27, 0x86, 0x20, 0xce, 0xb5, 0x4c, 0x65, 0xb2,
	0x1d, 0x65, 0x4a, 0x6a, 0x89, 0x3d, 0xbd, 0xcd, 0xc4, 0x7a, 0xe8, 0xcd, 0x65, 0xac, 0x7e, 0x96,
	0xf7, 0x86, 0x7e, 0xa6, 0xe4, 0x5f, 0xb1, 0xd0, 0x7c, 0x09, 0x2a, 0x5f, 0x09, 0xfe, 0xdc, 0x5f,
	0xfc, 0x89, 0xd3, 0xdf, 0x7c, 0xf5, 0xee, 0x1f, 0x40, 0x70, 0xcd, 0xbe, 0xeb, 0x85, 0x5e, 0xca,
	0x14, 0x3f, 0xc2, 0x51, 0xa6, 0x64, 0x76, 0x53, 0xe8, 0xc2, 0xbd, 0xb3, 0xbd, 0x73, 0xef, 0xe2,
	0x74, 0x44, 0xdf, 0x31, 0x9a, 0x29, 0x99, 0xc9, 0x75, 0xbc, 0xa2, 0xb3, 0x69, 0x27, 0xaa, 0x07,
	0xf1, 0x0a, 0xfa, 0x6a, 0x73, 0x3f, 0x33, 0x60, 0x97, 0xc0, 0x21, 0x83, 0x91, 0xd8, 0xc8, 0x7b,
	0xe1, 0xe2, 0x16, 0x81, 0x57, 0xe0, 0x6f, 0xa4, 0x16, 0xb5, 0xe2, 0x11, 0x29, 0x42, 0x56, 0xdc,
	0xf1, 0x59, 0x53, 0x60, 0x03, 0xf8, 0x15, 0x02, 0x9d, 0x2c, 0x53, 0x5d, 0x2b, 0x1e, 0x93, 0xe2,
	0x35, 0x2b, 0x6e, 0x85, 0x4a, 0x96, 0x69, 0xbc, 0xeb, 0x71, 0x30, 0xfc, 0x0c, 0x5e, 0xb1, 0xd9,
	0xac, 0x8c, 0x18, 0xf6, 0xc8, 0xf2, 0xcc, 0x89, 0xc0, 0xa7, 0xd3, 0x4e, 0xd4, 0x1c, 0xc6, 0x09,
	0x04, 0xbc, 0x56, 0x85, 0xef, 0x13, 0xfe, 0xaa, 0x35, 0x45, 0x2d, 0x71, 0x28, 0x9c, 0xc0, 0xa0,
	0xda, 0xae, 0x12, 0x1d, 0x58, 0x4d, 0x9b, 0x41, 0x6a, 0x8d, 0x0b, 0xe1, 0x0f, 0xc0, 0x2c, 0x9f,
	0xdf, 0x39, 0xaa, 0x43, 0x2b, 0xcc, 0xac, 0x1e, 0xb0, 0x6d, 0x2d, 0x28, 0x7e, 0x87, 0x63, 0x93,
	0xab, 0xd2, 0x1d, 0x91, 0xee, 0xed, 0x43, 0x9d, 0x6b, 0xe1, 0x0e, 0x8a, 0xef, 0xe1, 0xb0, 0xc8,
	0x17, 0xe5, 0x2b, 0x11, 0x02, 0x69, 0x4e, 0x9c, 0xd0, 0xc5, 0xd1, 0xb4, 0x13, 0x99, 0x31, 0xbc,
	0x04, 0x8f, 0x63, 0x11, 0xe5, 0x11, 0xf5, 0xa2, 0xb5, 0x2f, 0xb3, 0xcd, 0x79, 0xbc, 0x84, 0x7e,
	0x15, 0x89, 0xf8, 0x3e, 0xf1, 0xcf, 0x5b, 0xb2, 0x32, 0x6d, 0x8d, 0xe3, 0x17, 0xf0, 0xcd, 0x12,
	0xc4, 0xfb, 0xd6, 0xf3, 0xdd, 0x59, 0x9e, 0x25, 0x36, 0x54, 0xac, 0xad, 0x55, 0x9c, 0xae, 0x7f,
	0x09, 0x15, 0x06, 0xd6, 0xda, 0xb7, 0x7c, 0x7b, 0x92, 0xa7, 0xc5, 0x6f, 0xd3, 0x8c, 0xe1, 0x05,
	0x78, 0x0b, 0x99, 0x24, 0xa2, 0xb4, 0x84, 0x03, 0xa2, 0x02, 0xa6, 0xc6, 0xe5, 0x49, 0xb1, 0x6b,
	0x63, 0x08, 0x3f, 0x01, 0x14, 0xd9, 0xc6, 0xf4, 0x9f, 0x0f, 0x8f, 0x09, 0x79, 0xea, 0xf4, 0x2d,
	0x0f, 0xa7, 0x9d, 0xa8, 0x31, 0x8a, 0x63, 0xf0, 0xb9, 0x19, 0xb3, 0x4f, 0x88, 0x7d, 0xd9, 0x5a,
	0xd9, 0x18, 0x6c, 0x06, 0xc7, 0x10, 0x54, 0xe9, 0xd8, 0x82, 0xd6, 0xb3, 0x6a, 0xb6, 0x36, 0x0e,
	0x07, 0xc1, 0x6f, 0x30, 0x30, 0xe9, 0xd8, 0x72, 0x42, 0x96, 0x37, 0x0f, 0x15, 0x37, 0x2a, 0x17,
	0xc4, 0x00, 0xba, 0x7a, 0x1b, 0x9e, 0x9e, 0xed, 0x9d, 0xf7, 0xa2, 0xae, 0xde, 0xde, 0x1c, 0x40,
	0x6f, 0x13, 0xaf, 0x72, 0x31, 0xdf, 0xa7, 0xd7, 0xe1, 0x87, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x43, 0xb3, 0xd4, 0x92, 0x5d, 0x05, 0x00, 0x00,
}
