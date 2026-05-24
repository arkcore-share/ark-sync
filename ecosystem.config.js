module.exports = {
  apps: [{
    name: 'arksync',
    script: '/opt/ark-sync/bin/arksync',
    args: 'serve',
    env: {
      STGUIADDRESS: '0.0.0.0:8384',
      ARKSYNC_SKIP_PARENT_CHECK: '1',
    },
  }],
}
