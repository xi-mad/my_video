import request from '../request'

export const getConfig = () => {
    return request({
        url: '/api/config/',
        method: 'get',
    })
}