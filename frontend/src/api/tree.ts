import request from '../request'

export const createTreeNode = (data: any) => {
    return request({
        url: '/api/tree/create',
        method: 'post',
        data
    })
}

export const updateTreeNode = (data: any) => {
    return request({
        url: '/api/tree/update',
        method: 'put',
        data
    })
}

export const listTree = (data: any) => {
    return request({
        url: '/api/tree/list',
        method: 'get',
        data
    })
}



export const deleteTreeNode = (data: any) => {
    return request({
        url: '/api/tree/delete',
        method: 'delete',
        data
    })
}

export const optionsTree = () => {
    return request({
        url: '/api/tree/options',
        method: 'get',
    })
}