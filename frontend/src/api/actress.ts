import request from '../request'

export const createActress = (data: any) => {
    return request({
        url: '/api/actress/create',
        method: 'post',
        data
    })
}

export const updateActress = (data: any) => {
    return request({
        url: '/api/actress/update',
        method: 'put',
        data
    })
}

export const listActress = (data: any) => {
    const searchParams = new URLSearchParams(data);
    return request({
        url: '/api/actress/list?' + searchParams,
        method: 'get',
        data
    })
}

export const deleteActress = (data: any) => {
    return request({
        url: '/api/actress/delete',
        method: 'delete',
        data
    })
}

export const optionsActress = () => {
    return request({
        url: '/api/actress/options',
        method: 'get',
    })
}