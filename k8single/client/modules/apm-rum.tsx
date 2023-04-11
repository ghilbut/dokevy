import { init as initApm } from '@elastic/apm-rum'

const apm = initApm({
  serviceName: 'k8single-client-react',
  serviceVersion: 'v0.0',
})

export default apm;
