package gclient

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 客户端向服务端发 Hello 以确认服务端是否存活
type Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hello string `protobuf:"bytes,1,opt,name=hello,proto3" json:"hello,omitempty"`
}

func (x *Ping) Reset() {
	*x = Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_instago_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ping) ProtoMessage() {}

func (x *Ping) ProtoReflect() protoreflect.Message {
	mi := &file_instago_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ping.ProtoReflect.Descriptor instead.
func (*Ping) Descriptor() ([]byte, []int) {
	return file_instago_proto_rawDescGZIP(), []int{0}
}

func (x *Ping) GetHello() string {
	if x != nil {
		return x.Hello
	}
	return ""
}

// 服务端存活则回复
type Pong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Alive   bool   `protobuf:"varint,1,opt,name=alive,proto3" json:"alive,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Pong) Reset() {
	*x = Pong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_instago_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pong) ProtoMessage() {}

func (x *Pong) ProtoReflect() protoreflect.Message {
	mi := &file_instago_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pong.ProtoReflect.Descriptor instead.
func (*Pong) Descriptor() ([]byte, []int) {
	return file_instago_proto_rawDescGZIP(), []int{1}
}

func (x *Pong) GetAlive() bool {
	if x != nil {
		return x.Alive
	}
	return false
}

func (x *Pong) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// 客户端发送一个 Instagram 的分享地址给服务端
type ShareURL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *ShareURL) Reset() {
	*x = ShareURL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_instago_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShareURL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShareURL) ProtoMessage() {}

func (x *ShareURL) ProtoReflect() protoreflect.Message {
	mi := &file_instago_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShareURL.ProtoReflect.Descriptor instead.
func (*ShareURL) Descriptor() ([]byte, []int) {
	return file_instago_proto_rawDescGZIP(), []int{2}
}

func (x *ShareURL) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

// 服务端回复所有媒体地址给客户端
type MediaURLs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls []string `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *MediaURLs) Reset() {
	*x = MediaURLs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_instago_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaURLs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaURLs) ProtoMessage() {}

func (x *MediaURLs) ProtoReflect() protoreflect.Message {
	mi := &file_instago_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaURLs.ProtoReflect.Descriptor instead.
func (*MediaURLs) Descriptor() ([]byte, []int) {
	return file_instago_proto_rawDescGZIP(), []int{3}
}

func (x *MediaURLs) GetUrls() []string {
	if x != nil {
		return x.Urls
	}
	return nil
}

var File_instago_proto protoreflect.FileDescriptor

var file_instago_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x22, 0x1c, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x68,
	0x65, 0x6c, 0x6c, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x68, 0x65, 0x6c, 0x6c,
	0x6f, 0x22, 0x36, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69,
	0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x76, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x1c, 0x0a, 0x08, 0x53, 0x68, 0x61,
	0x72, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x1f, 0x0a, 0x09, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x55, 0x52, 0x4c, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x32, 0x5c, 0x0a, 0x0a, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x23, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x08, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x1a,
	0x08, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x08, 0x47,
	0x65, 0x74, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x0c, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x68, 0x61,
	0x72, 0x65, 0x55, 0x52, 0x4c, 0x1a, 0x0d, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x55, 0x52, 0x4c, 0x73, 0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d, 0x6c, 0x69, 0x62, 0x73, 0x2f, 0x3b,
	0x69, 0x6e, 0x73, 0x74, 0x61, 0x67, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_instago_proto_rawDescOnce sync.Once
	file_instago_proto_rawDescData = file_instago_proto_rawDesc
)

func file_instago_proto_rawDescGZIP() []byte {
	file_instago_proto_rawDescOnce.Do(func() {
		file_instago_proto_rawDescData = protoimpl.X.CompressGZIP(file_instago_proto_rawDescData)
	})
	return file_instago_proto_rawDescData
}

var file_instago_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_instago_proto_goTypes = []interface{}{
	(*Ping)(nil),      // 0: pb.Ping
	(*Pong)(nil),      // 1: pb.Pong
	(*ShareURL)(nil),  // 2: pb.ShareURL
	(*MediaURLs)(nil), // 3: pb.MediaURLs
}
var file_instago_proto_depIdxs = []int32{
	0, // 0: pb.InstaMedia.ServerCheck:input_type -> pb.Ping
	2, // 1: pb.InstaMedia.GetMedia:input_type -> pb.ShareURL
	1, // 2: pb.InstaMedia.ServerCheck:output_type -> pb.Pong
	3, // 3: pb.InstaMedia.GetMedia:output_type -> pb.MediaURLs
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_instago_proto_init() }
func file_instago_proto_init() {
	if File_instago_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_instago_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ping); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_instago_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pong); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_instago_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShareURL); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_instago_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaURLs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_instago_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_instago_proto_goTypes,
		DependencyIndexes: file_instago_proto_depIdxs,
		MessageInfos:      file_instago_proto_msgTypes,
	}.Build()
	File_instago_proto = out.File
	file_instago_proto_rawDesc = nil
	file_instago_proto_goTypes = nil
	file_instago_proto_depIdxs = nil
}
