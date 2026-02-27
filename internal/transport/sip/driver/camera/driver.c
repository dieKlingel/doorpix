#include "driver.h"

int Register(factory_options options) {

    pj_pool_factory *pool_factory;
	pj_pool_t *pool;
	factory_t *factory;

	pool_factory = pjsua_get_pool_factory();
	pool = pj_pool_create(pool_factory, "driver%p", 512, 512, NULL);
	if (pool == NULL) {
		return 1;
	}

	factory = PJ_POOL_ZALLOC_T(pool, factory_t);

	factory->pool_factory = pool_factory;
	factory->pool = pool;
	factory->base.op = &factory_op;
	factory->options = options;

	return pjmedia_vid_register_factory(NULL, &factory->base);
}
