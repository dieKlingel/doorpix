#include "driver.h"

void Register() {
    pj_pool_factory *pool_factory;
	pj_pool_t *pool;
	factory_t *factory;

	pool_factory = pjsua_get_pool_factory();
	pool = pj_pool_create(pool_factory, "appvideofc%p", 512, 512, NULL);
	factory = PJ_POOL_ZALLOC_T(pool, factory_t);

	factory->pool_factory = pjsua_get_pool_factory();
	factory->pool = pool;
	factory->base.op = &factory_op;

	pjmedia_vid_register_factory(NULL, &factory->base);
}
