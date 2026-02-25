#include "stream.h"

#include "driver.h"

pjmedia_vid_dev_stream_op stream_op = {
	.get_param = &stream_get_param,
	.get_cap = &stream_get_cap,
	.set_cap = &stream_set_cap,
	.start = &stream_start,
	.get_frame = &stream_get_frame,
	.put_frame = &stream_put_frame,
	.stop = &stream_stop,
	.destroy = &stream_destroy,
};

pj_status_t stream_get_frame(pjmedia_vid_dev_stream *stream, pjmedia_frame *frame)
{
	int success = go_stream_get_frame(stream, frame);
    if (success != PJ_SUCCESS) {
        return success;
    }

    frame->type = PJMEDIA_FRAME_TYPE_VIDEO;
	return PJ_SUCCESS;
}

pj_status_t stream_destroy(pjmedia_vid_dev_stream *pj_stream)
{
	stream_t *stream = (stream_t *)pj_stream;
	pj_pool_safe_release(&stream->pool);

	return PJ_SUCCESS;
}

pj_status_t stream_start(pjmedia_vid_dev_stream *stream)
{
    printf("stream_start\n");
	return go_stream_start(stream);
}

pj_status_t stream_stop(pjmedia_vid_dev_stream *stream)
{
	return go_stream_stop(stream);
}

pj_status_t stream_get_cap(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_cap cap, void *value)
{
	printf("stream_get_cap\n");
	return PJ_SUCCESS;
}

pj_status_t stream_set_cap(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_cap cap, const void *value)
{
	printf("stream_set_cap\n");
	return PJ_SUCCESS;
}

pj_status_t stream_get_param(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_param *param)
{
	printf("stream_get_param\n");
	return PJ_SUCCESS;
}

pj_status_t stream_put_frame(pjmedia_vid_dev_stream *stream, const pjmedia_frame *frame)
{
	return PJ_SUCCESS;
}
