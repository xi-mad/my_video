import request from '../request'

export const createTag = (data: any) => {
    return request({
        url: '/api/tag/create',
        method: 'post',
        data
    })
}

export const listTag = (data: any) => {
    const searchParams = new URLSearchParams(data);
    return request({
        url: '/api/tag/list?' + searchParams,
        method: 'get',
        data
    })
}

export const deleteTag = (data: any) => {
    return request({
        url: '/api/tag/delete',
        method: 'delete',
        data
    })
}

export const optionsTag = () => {
    return request({
        url: '/api/tag/options',
        method: 'get',
    })
}