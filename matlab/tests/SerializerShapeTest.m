classdef SerializerShapeTest < matlab.unittest.TestCase
    methods (Test)

        function testFixedVectorShape(testCase)
            fvs = yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 42);
            testCase.verifyEqual(fvs.getShape(), [1, 42]);

            fvs = yardl.binary.FixedVectorSerializer(fvs, 11);
            testCase.verifyEqual(fvs.getShape(), [42, 11]);

            fvs = yardl.binary.FixedVectorSerializer(yardl.binary.VectorSerializer(yardl.binary.Int32Serializer), 13);
            testCase.verifyEqual(fvs.getShape(), [1, 13]);

            fvs = yardl.binary.FixedVectorSerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3, 4]), 5);
            testCase.verifyEqual(fvs.getShape(), [3, 4, 5]);

            fvs = yardl.binary.FixedVectorSerializer(yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2), 6);
            testCase.verifyEqual(fvs.getShape(), [1, 6]);

            fvs = yardl.binary.FixedVectorSerializer(yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer), 7);
            testCase.verifyEqual(fvs.getShape(), [1, 7]);
        end

        function testVectorShape(testCase)
            vs = yardl.binary.VectorSerializer(yardl.binary.Int32Serializer);
            testCase.verifyEqual(vs.getShape(), []);

            vs = yardl.binary.VectorSerializer(vs);
            testCase.verifyEqual(vs.getShape(), []);

            vs = yardl.binary.VectorSerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 42));
            testCase.verifyEqual(vs.getShape(), []);

            fvs = yardl.binary.VectorSerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3, 4]));
            testCase.verifyEqual(fvs.getShape(), []);

            fvs = yardl.binary.VectorSerializer(yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2));
            testCase.verifyEqual(fvs.getShape(), []);

            fvs = yardl.binary.VectorSerializer(yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer));
            testCase.verifyEqual(fvs.getShape(), []);
        end

        function testFixedNDArrayShape(testCase)
            fas = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3, 4]);
            testCase.verifyEqual(fas.getShape(), [3, 4]);

            fas = yardl.binary.FixedNDArraySerializer(fas, [5, 6]);
            testCase.verifyEqual(fas.getShape(), [3, 4, 5, 6]);

            fas = yardl.binary.FixedNDArraySerializer(yardl.binary.VectorSerializer(yardl.binary.Int32Serializer), [7, 8]);
            testCase.verifyEqual(fas.getShape(), [7, 8]);

            fas = yardl.binary.FixedNDArraySerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 42), [7, 8]);
            testCase.verifyEqual(fas.getShape(), [42, 7, 8]);

            fas = yardl.binary.FixedNDArraySerializer(yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2), [7, 8]);
            testCase.verifyEqual(fas.getShape(), [7, 8]);

            fas = yardl.binary.FixedNDArraySerializer(yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer), [7, 8]);
            testCase.verifyEqual(fas.getShape(), [7, 8]);
        end

        function testNDArrayShape(testCase)
            nas = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
            testCase.verifyEqual(nas.getShape(), []);

            nas = yardl.binary.NDArraySerializer(nas, 3);
            testCase.verifyEqual(nas.getShape(), []);

            nas = yardl.binary.NDArraySerializer(yardl.binary.VectorSerializer(yardl.binary.Int32Serializer), 2);
            testCase.verifyEqual(nas.getShape(), []);

            nas = yardl.binary.NDArraySerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 42), 2);
            testCase.verifyEqual(nas.getShape(), []);

            nas = yardl.binary.NDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [5, 6]), 2);
            testCase.verifyEqual(nas.getShape(), []);

            nas = yardl.binary.NDArraySerializer(yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer), 2);
            testCase.verifyEqual(nas.getShape(), []);
        end

        function testDynamicNDArrayShape(testCase)
            das = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
            testCase.verifyEqual(das.getShape(), []);

            das = yardl.binary.DynamicNDArraySerializer(das);
            testCase.verifyEqual(das.getShape(), []);

            das = yardl.binary.DynamicNDArraySerializer(yardl.binary.VectorSerializer(yardl.binary.Int32Serializer));
            testCase.verifyEqual(das.getShape(), []);

            das = yardl.binary.DynamicNDArraySerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 42));
            testCase.verifyEqual(das.getShape(), []);

            das = yardl.binary.DynamicNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [5, 6]));
            testCase.verifyEqual(das.getShape(), []);

            das = yardl.binary.DynamicNDArraySerializer(yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2));
            testCase.verifyEqual(das.getShape(), []);
        end
    end
end
