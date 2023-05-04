import request from '../request'

export const createCollection = (data: any) => {
    return request({
        url: '/api/collection/',
        method: 'post',
        data
    })
}

export const listCollection = (data: any) => {
    const searchParams = new URLSearchParams(data);
    return request({
        url: '/api/collection?' + searchParams,
        method: 'get',
    })
}

export const getCollection = (id: any) => {
    return request({
        url: '/api/collection/get?id=' + id,
        method: 'get',
    })
}

export const updateCollection = (data: any) => {
    return request({
        url: '/api/collection/',
        method: 'put',
        data
    })
}

export const deleteCollection = (data: any) => {
    return request({
        url: '/api/collection/',
        method: 'delete',
        data
    })
}

export const optionsCollection = () => {
    return request({
        url: '/api/collection/options',
        method: 'get'
    })
}

export const associateCollection = (data: any) => {
    return request({
        url: '/api/collection/associate',
        method: 'post',
        data
    })
}

export const disassociateCollection = (data: any) => {
    return request({
        url: '/api/collection/disassociate',
        method: 'post',
        data
    })
}

export const detailCollection = (id: any) => {
    return request({
        url: '/api/collection/detail?id=' + id,
        method: 'get',
    })
}