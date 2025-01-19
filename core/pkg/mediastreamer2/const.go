package mediastreamer2

// #include <mediastreamer2/allfilters.h>
// #include <mediastreamer2/mscommon.h>
// #include <mediastreamer2/msjpegwriter.h>
// #include <mediastreamer2/msticker.h>
// #include <mediastreamer2/msvideo.h>
//#include <mediastreamer2/mswebcam.h>
import "C"

type FilterMethod = int

type FilterId = C.MSFilterId

type FilterInterfaceId = C.MSFilterInterfaceId

const MS_FILTER_GET_VIDEO_SIZE FilterMethod = C.MS_FILTER_GET_VIDEO_SIZE
const MS_FILTER_SET_VIDEO_SIZE FilterMethod = C.MS_FILTER_SET_VIDEO_SIZE
const MS_STATIC_IMAGE_ID FilterId = C.MS_STATIC_IMAGE_ID
const MS_FILTER_SET_FPS FilterMethod = C.MS_FILTER_SET_FPS
const MS_VIDEO_ENCODER_GET_CONFIGURATION FilterMethod = C.MS_VIDEO_ENCODER_GET_CONFIGURATION
const MS_VIDEO_ENCODER_SET_CONFIGURATION FilterMethod = C.MS_VIDEO_ENCODER_SET_CONFIGURATION
const MS_FILTER_GET_PIX_FMT FilterMethod = C.MS_FILTER_GET_PIX_FMT
const MS_MJPEG FilterId = C.MS_MJPEG
const MS_MJPEG_DEC_ID FilterId = C.MS_MJPEG_DEC_ID
const MS_PIX_CONV_ID FilterId = C.MS_PIX_CONV_ID
const MS_FILTER_SET_PIX_FMT FilterMethod = C.MS_FILTER_SET_PIX_FMT
const MS_JPEG_WRITER_ID FilterId = C.MS_JPEG_WRITER_ID
const MS_JPEG_WRITER_TAKE_SNAPSHOT FilterMethod = C.MS_JPEG_WRITER_TAKE_SNAPSHOT
const MSFilterVideoEncoderInterface FilterInterfaceId = C.MSFilterVideoEncoderInterface
