// Compatibility shim: start the CRA dev server on webpack-dev-server v5.
//
// react-scripts 5.0.1 configures the dev server with `onBeforeSetupMiddleware`
// and `onAfterSetupMiddleware`, both removed in webpack-dev-server v5. We are
// pinned to v5 because every webpack-dev-server security fix (GHSA-79cf-xcqc-c78w,
// GHSA-mx8g-39q3-5c79, GHSA-f5vj-f2hx-8m93, GHSA-m28w-2pqf-7qgj) landed only in
// the 5.x line -- 4.15.2 is the last 4.x release and stays vulnerable.
//
// CRA is unmaintained (5.0.1 is the final release), so instead of patching
// node_modules we translate the two old hooks into v5's `setupMiddlewares`
// before react-scripts hands the object to webpack-dev-server.
const path = require('path');

const configPath = require.resolve('react-scripts/config/webpackDevServer.config.js');
const original = require(configPath);

require.cache[configPath].exports = function patchedConfig(...args) {
  const config = original(...args);
  const before = config.onBeforeSetupMiddleware;
  const after = config.onAfterSetupMiddleware;
  delete config.onBeforeSetupMiddleware;
  delete config.onAfterSetupMiddleware;

  // v5 renamed `https` -> `server` and removed the standalone `http2` flag.
  if ('https' in config) {
    const https = config.https;
    delete config.https;
    if (https) {
      config.server =
        typeof https === 'object' ? { type: 'https', options: https } : 'https';
    }
  }
  delete config.http2;

  config.setupMiddlewares = (middlewares, devServer) => {
    // v5 drops `devServer.app`; the old hooks call `.use()` on it, so map those
    // registrations onto the middleware array in the same relative order.
    const prepended = [];
    const appended = [];
    const collector = sink => ({ use: (...handlers) => sink.push(...handlers) });

    if (before) before({ ...devServer, app: collector(prepended) });
    if (after) after({ ...devServer, app: collector(appended) });

    const result = [
      ...prepended.map(fn => ({ name: 'cra-before', middleware: fn })),
      ...middlewares,
    ];

    // In v4 the `after` hook ran before `connect-history-api-fallback`; in v5 the
    // fallback is pushed ahead of `webpack-dev-middleware`, so appending at the end
    // would put these behind a catch-all that answers every GET with index.html.
    // Splice them in just before the fallback to keep the v4 ordering.
    const fallbackAt = result.findIndex(
      m => m && m.name === 'connect-history-api-fallback'
    );
    const insertAt = fallbackAt === -1 ? result.length : fallbackAt;
    result.splice(
      insertAt,
      0,
      ...appended.map(fn => ({ name: 'cra-after', middleware: fn }))
    );

    return result;
  };

  return config;
};

require(path.resolve(__dirname, 'node_modules/react-scripts/scripts/start.js'));
