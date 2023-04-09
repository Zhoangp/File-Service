package https

import (
	"context"
	"fmt"
	"github.com/Zhoangp/File-Service/pb"
	"github.com/Zhoangp/File-Service/pkg/common"
	"github.com/Zhoangp/File-Service/pkg/upload"
)

type uploadHandler struct {
	s3Provider upload.UploadProvider
	pb.UnimplementedFileServiceServer
}

func NewUploadHandler(s3Provider upload.UploadProvider) *uploadHandler {
	return &uploadHandler{s3Provider: s3Provider}
}
func (hdl *uploadHandler) UploadFile(ctx context.Context, rq *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
		// Lấy width, height của image
		// w, h, err := getImageDemension(dataBytes)
		// if err != nil {
		// 	// File không phải là hình ảnh
		// 	panic(common.ErrBadRequest(err))
		// }
		url, err := hdl.s3Provider.UploadFile(
			ctx,
			rq.Content,
			fmt.Sprintf("%s/%s", rq.Folder, rq.FileName),
		)

		if err != nil {
			return &pb.UploadFileResponse{
				Error: HandleError(err),
			}, nil
		}
		return &pb.UploadFileResponse{
			Url: url,
		}, nil
}

func HandleError(err error) *pb.ErrorResponse {
	if errors, ok := err.(*common.AppError); ok {
		return &pb.ErrorResponse{
			Code:    int64(errors.StatusCode),
			Message: errors.Message,
		}
	}
	appErr := common.ErrInternal(err.(error))
	return &pb.ErrorResponse{
		Code:    int64(appErr.StatusCode),
		Message: appErr.Message,
	}
}
