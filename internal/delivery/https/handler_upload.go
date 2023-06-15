package https

import (
	"context"
	"fmt"
	"github.com/Zhoangp/File-Service/config"
	"github.com/Zhoangp/File-Service/pb"
	"github.com/Zhoangp/File-Service/pkg/common"
	"github.com/Zhoangp/File-Service/pkg/upload"
	"github.com/Zhoangp/File-Service/pkg/utils"
	"io"
	"math"
)

type uploadHandler struct {
	s3Provider upload.UploadProvider
	cf         *config.Config
	pb.UnimplementedFileServiceServer
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
func NewUploadHandler(s3Provider upload.UploadProvider, cf *config.Config) *uploadHandler {
	return &uploadHandler{s3Provider: s3Provider, cf: cf}
}
func (hdl *uploadHandler) UploadAvatar(ctx context.Context, rq *pb.UploadAvatarRequest) (*pb.UploadAvatarResponse, error) {
	if err := hdl.s3Provider.DeleteFile(rq.OldUrl); err != nil {
		fmt.Println(err)
		return &pb.UploadAvatarResponse{
			Error: HandleError(err),
		}, nil
	}
	url, err := hdl.s3Provider.UploadFile(
		ctx,
		rq.File.Content,
		fmt.Sprintf("%s/%s", rq.File.Folder, rq.File.FileName),
	)

	if err != nil {
		fmt.Println(err)
		return &pb.UploadAvatarResponse{
			Error: HandleError(err),
		}, nil
	}

	return &pb.UploadAvatarResponse{
		Url: url,
	}, nil
}

func (hdl *uploadHandler) UploadAsset(stream pb.FileService_UploadAssetServer) error {
	var file pb.File
	rq, err := stream.Recv()
	file.FileName = rq.File.FileName
	file.Size = rq.File.Size
	file.Folder = rq.File.Folder
	if err != nil {
		return err
	}
	for {
		rq, err = stream.Recv()
		if err == io.EOF {
			url, err := hdl.s3Provider.UploadFile(
				context.Background(),
				file.Content,
				fmt.Sprintf("%s/%s", file.Folder, file.FileName),
			)
			if err != nil {
				return err
			}
			var video pb.Video
			if file.Folder == "video" {
				v, err := utils.GetDurationVideo(hdl.cf, url)
				if err != nil {
					return err
				}
				video.Duration = int64(math.Ceil(v.Duration))
				video.Unit = "seconds"
			}
			if err = stream.Send(&pb.UploadAssetResponse{
				Url:   url,
				Video: &video,
			}); err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}

		file.Content = append(file.Content, rq.File.Content...)
	}

}
