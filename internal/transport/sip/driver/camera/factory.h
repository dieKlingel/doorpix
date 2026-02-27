#ifndef FACTORY_H
#define FACTORY_H

#include <pjlib.h>
#include <pjmedia/endpoint.h>
#include <pjmedia/vid_stream.h>
#include <pjmedia-videodev/videodev_imp.h>
#include <pjsua-lib/pjsua.h>

#include "stream.h"

pj_status_t factory_init(pjmedia_vid_dev_factory *f);
pj_status_t factory_destroy(pjmedia_vid_dev_factory *f);
unsigned int factory_get_dev_count(pjmedia_vid_dev_factory *f);
pj_status_t factory_get_dev_info(pjmedia_vid_dev_factory *f, unsigned int index, pjmedia_vid_dev_info *info);
pj_status_t factory_create_stream(pjmedia_vid_dev_factory *f, pjmedia_vid_dev_param *param, const pjmedia_vid_dev_cb *cb, void *user_data, pjmedia_vid_dev_stream **p_vid_strm);
pj_status_t factory_default_param(pj_pool_t *pool, pjmedia_vid_dev_factory *f, unsigned index, pjmedia_vid_dev_param *param);
pj_status_t factory_refresh(pjmedia_vid_dev_factory *f);

typedef struct factory_options {
	char *name;
	int name_length;
	char *driver_name;
	int driver_name_length;
	int width;
	int height;
	int framerate;
} factory_options;

typedef struct factory
{
	pjmedia_vid_dev_factory base;

	pj_pool_t *pool;
	pj_pool_factory *pool_factory;
	factory_options options;
} factory_t;

extern pjmedia_vid_dev_factory_op factory_op;

#endif // FACTORY_H
