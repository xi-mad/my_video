import request from '../request'

export const createObject = (data: any) => {
    return request({
        url: '/api/object/create',
        method: 'post',
        data
    })
}

export const updateObject = (data: any) => {
    return request({
        url: '/api/object/update',
        method: 'put',
        data
    })
}

export const deleteObject = (data: any) => {
    return request({
        url: '/api/object/delete',
        method: 'delete',
        data
    })
}

export const listObject = (data: any) => {
    return request({
        url: '/api/object/list',
        method: 'post',
        data
    })
}

export const scanObject = (data: any) => {
    return request({
        url: '/api/object/scan',
        method: 'post',
        data
    })
}

export const playPath = (path: any) => {
    return request({
        url: '/api/object/play?path=' + path,
        method: 'get'
    })
}

export const log = () => {
    return request({
        url: '/api/object/log',
        method: 'get'
    })
}

export const videoPath = (path: any) => {
    return request.getUri({
        url: '/api/object/video?path=' + path,
    })
}

export const viewinc = (id: any) => {
    return request({
        url: '/api/object/viewinc?id=' + id,
        method: 'get'
    })
}

export const randomPath = () => {
    return request({
        url: '/api/object/random',
        method: 'get'
    })
}

export const addTags = (data: any) => {
    return request({
        url: '/api/object/addTags',
        method: 'post',
        data
    })
}