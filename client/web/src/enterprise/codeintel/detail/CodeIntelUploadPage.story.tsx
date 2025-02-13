import { storiesOf } from '@storybook/react'
import React from 'react'
import { Observable, of } from 'rxjs'

import { LSIFIndexState } from '@sourcegraph/shared/src/graphql/schema'

import { LsifUploadConnectionFields, LsifUploadFields, LSIFUploadState } from '../../../graphql-operations'
import { EnterpriseWebStory } from '../../components/EnterpriseWebStory'

import { CodeIntelUploadPage } from './CodeIntelUploadPage'

const { add } = storiesOf('web/codeintel/detail/CodeIntelUploadPage', module)
    .addDecorator(story => <div className="p-3 container">{story()}</div>)
    .addParameters({
        chromatic: {
            viewports: [320, 576, 978, 1440],
        },
    })

add('Uploading', () => (
    <EnterpriseWebStory>
        {props => (
            <CodeIntelUploadPage
                {...props}
                fetchLsifUpload={fetchUpload({
                    state: LSIFUploadState.UPLOADING,
                    uploadedAt: '2020-06-15T12:20:30+00:00',
                    startedAt: null,
                    finishedAt: null,
                    failure: null,
                    placeInQueue: null,
                    associatedIndex: null,
                })}
                fetchLsifUploads={fetchEmptyDependencies}
                now={now}
            />
        )}
    </EnterpriseWebStory>
))

add('Queued', () => (
    <EnterpriseWebStory>
        {props => (
            <CodeIntelUploadPage
                {...props}
                fetchLsifUpload={fetchUpload({
                    state: LSIFUploadState.QUEUED,
                    uploadedAt: '2020-06-15T12:20:30+00:00',
                    startedAt: null,
                    finishedAt: null,
                    placeInQueue: 3,
                    failure: null,
                    associatedIndex: null,
                })}
                fetchLsifUploads={fetchEmptyDependencies}
                now={now}
            />
        )}
    </EnterpriseWebStory>
))

add('Processing', () => (
    <EnterpriseWebStory>
        {props => (
            <CodeIntelUploadPage
                {...props}
                fetchLsifUpload={fetchUpload({
                    state: LSIFUploadState.PROCESSING,
                    uploadedAt: '2020-06-15T12:20:30+00:00',
                    startedAt: '2020-06-15T12:25:30+00:00',
                    finishedAt: null,
                    failure: null,
                    placeInQueue: null,
                    associatedIndex: null,
                })}
                fetchLsifUploads={fetchEmptyDependencies}
                now={now}
            />
        )}
    </EnterpriseWebStory>
))

add('Completed', () => {
    const upload = fetchUpload({
        state: LSIFUploadState.COMPLETED,
        uploadedAt: '2020-06-15T12:20:30+00:00',
        startedAt: '2020-06-15T12:25:30+00:00',
        finishedAt: '2020-06-15T12:30:30+00:00',
        failure: null,
        placeInQueue: null,
        associatedIndex: null,
    })
    const uploads = fetchUploads({ id: '1' }, { id: '2' }, { id: '3' }, { id: '4' }, { id: '5' }, { id: '6' })

    return (
        <EnterpriseWebStory>
            {props => <CodeIntelUploadPage {...props} fetchLsifUpload={upload} fetchLsifUploads={uploads} now={now} />}
        </EnterpriseWebStory>
    )
})

add('Errored', () => (
    <EnterpriseWebStory>
        {props => (
            <CodeIntelUploadPage
                {...props}
                fetchLsifUpload={fetchUpload({
                    state: LSIFUploadState.ERRORED,
                    uploadedAt: '2020-06-15T12:20:30+00:00',
                    startedAt: '2020-06-15T12:25:30+00:00',
                    finishedAt: '2020-06-15T12:30:30+00:00',
                    failure:
                        'Upload failed to complete: dial tcp: lookup gitserver-8.gitserver on 10.165.0.10:53: no such host',
                    placeInQueue: null,
                    associatedIndex: null,
                })}
                fetchLsifUploads={fetchEmptyDependencies}
                now={now}
            />
        )}
    </EnterpriseWebStory>
))

add('Deleting', () => {
    const upload = fetchUpload({
        state: LSIFUploadState.DELETING,
        uploadedAt: '2020-06-15T12:20:30+00:00',
        startedAt: '2020-06-15T12:25:30+00:00',
        finishedAt: '2020-06-15T12:30:30+00:00',
        failure: null,
        placeInQueue: null,
        associatedIndex: null,
    })
    const uploads = fetchUploads({ id: '1' }, { id: '2' }, { id: '3' }, { id: '4' }, { id: '5' }, { id: '6' })

    return (
        <EnterpriseWebStory>
            {props => <CodeIntelUploadPage {...props} fetchLsifUpload={upload} fetchLsifUploads={uploads} now={now} />}
        </EnterpriseWebStory>
    )
})

add('Failed Upload', () => (
    <EnterpriseWebStory>
        {props => (
            <CodeIntelUploadPage
                {...props}
                fetchLsifUpload={fetchUpload({
                    state: LSIFUploadState.ERRORED,
                    uploadedAt: '2020-06-15T12:20:30+00:00',
                    startedAt: null,
                    finishedAt: '2020-06-15T12:30:30+00:00',
                    failure: 'Upload failed to complete: object store error:\n * XMinioStorageFull etc etc',
                    placeInQueue: null,
                    associatedIndex: null,
                })}
                fetchLsifUploads={fetchEmptyDependencies}
                now={now}
            />
        )}
    </EnterpriseWebStory>
))

add('Associated Index', () => (
    <EnterpriseWebStory>
        {props => (
            <CodeIntelUploadPage
                {...props}
                fetchLsifUpload={fetchUpload({
                    state: LSIFUploadState.PROCESSING,
                    uploadedAt: '2020-06-15T12:20:30+00:00',
                    startedAt: '2020-06-15T12:25:30+00:00',
                    finishedAt: null,
                    failure: null,
                    placeInQueue: null,
                    associatedIndex: {
                        id: '6789',
                        state: LSIFIndexState.COMPLETED,
                        queuedAt: '2020-06-15T12:15:10+00:00',
                        startedAt: '2020-06-15T12:20:20+00:00',
                        finishedAt: '2020-06-15T12:25:30+00:00',
                        placeInQueue: null,
                    },
                })}
                fetchLsifUploads={fetchEmptyDependencies}
                now={now}
            />
        )}
    </EnterpriseWebStory>
))

const fetchUpload = (
    upload: Pick<
        LsifUploadFields,
        'state' | 'uploadedAt' | 'startedAt' | 'finishedAt' | 'failure' | 'placeInQueue' | 'associatedIndex'
    >
): (() => Observable<LsifUploadFields>) => () =>
    of({
        __typename: 'LSIFUpload',
        id: '1234',
        projectRoot: {
            url: '',
            path: 'web/',
            repository: {
                url: '',
                name: 'github.com/sourcegraph/sourcegraph',
            },
            commit: {
                url: '',
                oid: '9ea5e9f0e0344f8197622df6b36faf48ccd02570',
                abbreviatedOID: '9ea5e9f',
            },
        },
        inputCommit: '9ea5e9f0e0344f8197622df6b36faf48ccd02570',
        inputRoot: 'web/',
        inputIndexer: 'lsif-tsc',
        isLatestForRepo: false,
        ...upload,
    })

const fetchUploads = (
    ...uploads: Omit<
        LsifUploadFields,
        | '__typename'
        | 'projectRoot'
        | 'inputCommit'
        | 'inputRoot'
        | 'inputIndexer'
        | 'isLatestForRepo'
        | 'state'
        | 'uploadedAt'
        | 'startedAt'
        | 'finishedAt'
        | 'failure'
        | 'placeInQueue'
        | 'associatedIndex'
    >[]
): (() => Observable<LsifUploadConnectionFields>) => () =>
    of({
        nodes: uploads.map(upload => ({
            __typename: 'LSIFUpload',
            projectRoot: {
                url: '',
                path: 'web/',
                repository: {
                    url: '',
                    name: 'github.com/sourcegraph/sourcegraph',
                },
                commit: {
                    url: '',
                    oid: '9ea5e9f0e0344f8197622df6b36faf48ccd02570',
                    abbreviatedOID: '9ea5e9f',
                },
            },
            inputCommit: '9ea5e9f0e0344f8197622df6b36faf48ccd02570',
            inputRoot: 'web/',
            inputIndexer: 'lsif-tsc',
            isLatestForRepo: false,
            state: LSIFUploadState.COMPLETED,
            uploadedAt: '2020-06-15T12:20:30+00:00',
            startedAt: '2020-06-15T12:25:30+00:00',
            finishedAt: '2020-06-15T12:30:30+00:00',
            failure: null,
            placeInQueue: null,
            associatedIndex: null,
            ...upload,
        })),
        totalCount: uploads.length > 0 ? uploads.length + 5 : 0,
        pageInfo: {
            __typename: 'PageInfo',
            endCursor: uploads.length > 0 ? 'fakenextpage' : null,
            hasNextPage: uploads.length > 0,
        },
    })

const fetchEmptyDependencies = fetchUploads()

const now = () => new Date('2020-06-15T15:25:00+00:00')
